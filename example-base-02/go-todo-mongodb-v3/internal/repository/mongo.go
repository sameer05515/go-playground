package repository

import (
 "context"
 "errors"
 "strings"
 "time"
 "go-todo-mongodb-v3/internal/model"
 "go.mongodb.org/mongo-driver/v2/bson"
 "go.mongodb.org/mongo-driver/v2/mongo"
 "go.mongodb.org/mongo-driver/v2/mongo/options"
)
var ErrNotFound=errors.New("todo not found")
type Repository struct { Collection *mongo.Collection }
func New(db *mongo.Database)*Repository{return &Repository{Collection:db.Collection("todos")}}
func (r *Repository) List(ctx context.Context, completed *bool, search string)([]model.Todo,error){
 filter:=bson.M{}
 if completed!=nil {filter["completed"]=*completed}
 if strings.TrimSpace(search)!="" { filter["$or"]=bson.A{bson.M{"title":bson.M{"$regex":search,"$options":"i"}},bson.M{"description":bson.M{"$regex":search,"$options":"i"}}} }
 cursor,err:=r.Collection.Find(ctx,filter,options.Find().SetSort(bson.D{{Key:"createdAt",Value:-1}}));if err!=nil{return nil,err};defer cursor.Close(ctx)
 todos:=[]model.Todo{};if err=cursor.All(ctx,&todos);err!=nil{return nil,err};return todos,nil
}
func (r *Repository) Get(ctx context.Context,id bson.ObjectID)(model.Todo,error){var t model.Todo;err:=r.Collection.FindOne(ctx,bson.M{"_id":id}).Decode(&t);if errors.Is(err,mongo.ErrNoDocuments){return t,ErrNotFound};return t,err}
func (r *Repository) Create(ctx context.Context,t model.Todo)(model.Todo,error){t.ID=bson.NewObjectID();_,err:=r.Collection.InsertOne(ctx,t);return t,err}
func (r *Repository) Update(ctx context.Context,id bson.ObjectID,changes bson.M)(model.Todo,error){
 opts:=options.FindOneAndUpdate().SetReturnDocument(options.After)
 var t model.Todo;err:=r.Collection.FindOneAndUpdate(ctx,bson.M{"_id":id},bson.M{"$set":changes},opts).Decode(&t)
 if errors.Is(err,mongo.ErrNoDocuments){return t,ErrNotFound};return t,err
}
func (r *Repository) Delete(ctx context.Context,id bson.ObjectID)error{res,err:=r.Collection.DeleteOne(ctx,bson.M{"_id":id});if err!=nil{return err};if res.DeletedCount==0{return ErrNotFound};return nil}
func WithTimeout()(context.Context,context.CancelFunc){return context.WithTimeout(context.Background(),10*time.Second)}
