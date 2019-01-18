package main

import (
  "github.com/mongodb/mongo-go-driver/mongo"
  "github.com/mongodb/mongo-go-driver/bson"
  "github.com/satori/go.uuid"
  "log"
  "context"
  "time"
)

func connect() (*mongo.Client, error) {
  client, err := mongo.NewClient("mongodb://localhost:27017")
  if err != nil { return nil, err }
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  err = client.Connect(ctx)
  if err != nil { return nil, err }
  return client, nil
}

func setup_db() error {
  client, err := connect()
  if err != nil { return err }
  ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
  collection := client.Database("default").Collection("products")
  _, err = collection.DeleteMany(ctx, bson.M{})
  if err != nil { return err }
  var products = []interface{}{
    &Product{
      Uuid: uuid.Must(uuid.NewV4()).String(),
      Title: "Plastic Bag",
      Price: "1",
      Inventorycount: "100",
    },
    &Product{
      Uuid: uuid.Must(uuid.NewV4()).String(),
      Title: "Orange Peels",
      Price: "5",
      Inventorycount: "8",
    },
  }
  _, err = collection.InsertMany(ctx, products)
  return err
}

func get_record(collection_name string, record_uuid string) (bson.M, error) {
  filter := bson.D{{"uuid", record_uuid}}
  client, err := connect()
  if err != nil { return nil, err }
  collection := client.Database("default").Collection(collection_name)
  ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
  var result bson.M
  err = collection.FindOne(ctx, filter).Decode(&result)
  if err != nil { return nil, err }
  return result, nil
}

func get_records(collection_name string) ([]bson.M, error) {
  client, err := connect()
  if err != nil { return nil, err }
  var items []bson.M
  collection := client.Database("default").Collection(collection_name)
  ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
  cur, err := collection.Find(ctx, bson.M{})
  if err != nil { return nil, err }
  log.Print("iterating through items")
  defer cur.Close(ctx)
  for cur.Next(ctx) {
    log.Print("adding item to list")
    var result bson.M
    err = cur.Decode(&result)
    log.Print(result)
    if err != nil { return nil, err }
    items = append(items, result)
  }
  if err = cur.Err(); err != nil {
    return nil, err
  }
  return items, nil
}
