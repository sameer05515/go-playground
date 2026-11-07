package service

import (
 "context"
 "errors"
 "strings"
 "time"
 "go-todo-mongodb-v3/internal/model"
 "go-todo-mongodb-v3/internal/repository"
 "go.mongodb.org/mongo-driver/v2/bson"
)
var ErrTitle=errors.New("title is required")
type Service struct{ Repo *repository.Repository }
func New(r *repository.Repository)*Service{return &Service{Repo:r}}
func(s *Service) List(ctx context.Context,done *bool,q string)([]model.Todo,error){return s.Repo.List(ctx,done,q)}
func(s *Service) Get(ctx context.Context,id bson.ObjectID)(model.Todo,error){return s.Repo.Get(ctx,id)}
func(s *Service) Create(ctx context.Context,in model.Input)(model.Todo,error){title:=strings.TrimSpace(in.Title);if title==""{return model.Todo{},ErrTitle};now:=time.Now().UTC();return s.Repo.Create(ctx,model.Todo{Title:title,Description:strings.TrimSpace(in.Description),Completed:in.Completed,CreatedAt:now,UpdatedAt:now})}
func(s *Service) Update(ctx context.Context,id bson.ObjectID,in model.Input)(model.Todo,error){title:=strings.TrimSpace(in.Title);if title==""{return model.Todo{},ErrTitle};return s.Repo.Update(ctx,id,bson.M{"title":title,"description":strings.TrimSpace(in.Description),"completed":in.Completed,"updatedAt":time.Now().UTC()})}
func(s *Service) Patch(ctx context.Context,id bson.ObjectID,in model.Patch)(model.Todo,error){changes:=bson.M{"updatedAt":time.Now().UTC()};if in.Title!=nil{title:=strings.TrimSpace(*in.Title);if title==""{return model.Todo{},ErrTitle};changes["title"]=title};if in.Description!=nil{changes["description"]=strings.TrimSpace(*in.Description)};if in.Completed!=nil{changes["completed"]=*in.Completed};return s.Repo.Update(ctx,id,changes)}
func(s *Service) Delete(ctx context.Context,id bson.ObjectID)error{return s.Repo.Delete(ctx,id)}
