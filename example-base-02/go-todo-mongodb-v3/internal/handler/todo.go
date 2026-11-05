package handler

import (
 "encoding/json"
 "errors"
 "io"
 "net/http"
 "strconv"
 "strings"
 "go-todo-mongodb-v3/internal/model"
 "go-todo-mongodb-v3/internal/repository"
 "go-todo-mongodb-v3/internal/service"
 "go.mongodb.org/mongo-driver/v2/bson"
)
type Handler struct{ Service *service.Service }
func New(s *service.Service)*Handler{return &Handler{Service:s}}
func(h *Handler) Register(m *http.ServeMux){m.HandleFunc("GET /api/todos",h.list);m.HandleFunc("GET /api/todos/{id}",h.get);m.HandleFunc("POST /api/todos",h.create);m.HandleFunc("PUT /api/todos/{id}",h.update);m.HandleFunc("PATCH /api/todos/{id}",h.patch);m.HandleFunc("DELETE /api/todos/{id}",h.delete)}
func respond(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
func fail(w http.ResponseWriter,status int,msg string){respond(w,status,map[string]string{"error":msg})}
func decode(w http.ResponseWriter,r *http.Request,v any)bool{r.Body=http.MaxBytesReader(w,r.Body,1<<20);d:=json.NewDecoder(r.Body);d.DisallowUnknownFields();if err:=d.Decode(v);err!=nil{fail(w,400,"invalid JSON: "+err.Error());return false};if err:=d.Decode(&struct{}{});err!=io.EOF{fail(w,400,"expected exactly one JSON object");return false};return true}
func id(w http.ResponseWriter,r *http.Request)(bson.ObjectID,bool){oid,err:=bson.ObjectIDFromHex(r.PathValue("id"));if err!=nil{fail(w,400,"invalid MongoDB ObjectID");return bson.ObjectID{},false};return oid,true}
func handleError(w http.ResponseWriter,err error){switch{case errors.Is(err,repository.ErrNotFound):fail(w,404,"todo not found");case errors.Is(err,service.ErrTitle):fail(w,400,"title is required");default:fail(w,500,"internal server error")}}
func(h *Handler) list(w http.ResponseWriter,r *http.Request){var completed *bool;if raw:=r.URL.Query().Get("completed");raw!=""{v,err:=strconv.ParseBool(raw);if err!=nil{fail(w,400,"completed must be true or false");return};completed=&v};todos,err:=h.Service.List(r.Context(),completed,strings.TrimSpace(r.URL.Query().Get("q")));if err!=nil{handleError(w,err);return};respond(w,200,todos)}
func(h *Handler) get(w http.ResponseWriter,r *http.Request){oid,ok:=id(w,r);if !ok{return};todo,err:=h.Service.Get(r.Context(),oid);if err!=nil{handleError(w,err);return};respond(w,200,todo)}
func(h *Handler) create(w http.ResponseWriter,r *http.Request){var in model.Input;if !decode(w,r,&in){return};todo,err:=h.Service.Create(r.Context(),in);if err!=nil{handleError(w,err);return};respond(w,201,todo)}
func(h *Handler) update(w http.ResponseWriter,r *http.Request){oid,ok:=id(w,r);if !ok{return};var in model.Input;if !decode(w,r,&in){return};todo,err:=h.Service.Update(r.Context(),oid,in);if err!=nil{handleError(w,err);return};respond(w,200,todo)}
func(h *Handler) patch(w http.ResponseWriter,r *http.Request){oid,ok:=id(w,r);if !ok{return};var in model.Patch;if !decode(w,r,&in){return};todo,err:=h.Service.Patch(r.Context(),oid,in);if err!=nil{handleError(w,err);return};respond(w,200,todo)}
func(h *Handler) delete(w http.ResponseWriter,r *http.Request){oid,ok:=id(w,r);if !ok{return};if err:=h.Service.Delete(r.Context(),oid);err!=nil{handleError(w,err);return};w.WriteHeader(204)}
