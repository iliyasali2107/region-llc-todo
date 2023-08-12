package api

// func RegisterRoutes(r *gin.Engine, c *config.Config) *ServiceClient {
// 	svc := &ServiceClient{
// 		Client: InitServiceClient(c),
// 	}

// 	r.POST("/todo", svc.Create)
// 	r.PUT("/todo/tasks/:id", svc.Update)
// 	r.DELETE("/todo/tasks/:id", svc.Delete)
// 	r.PUT("/todo/tasks/:id/done", svc.UpdateAsDone)
// 	r.GET("todo/tasks", svc.ListTodos)

// 	return svc
// }

// func (svc *ServiceClient) Create(ctx *gin.Context) {
// 	handlers.Create(ctx, svc.Client)
// }

// func (svc *ServiceClient) Update(ctx *gin.Context) {
// 	handlers.Update(ctx, svc.Client)
// }

// func (svc *ServiceClient) Delete(ctx *gin.Context) {
// 	handlers.Delete(ctx, svc.Client)
// }

// func (svc *ServiceClient) UpdateAsDone(ctx *gin.Context) {
// 	handlers.UpdateAsDone(ctx, svc.Client)
// }

// func (svc *ServiceClient) ListTodos(ctx *gin.Context) {
// 	handlers.ListTodos(ctx, svc.Client)
// }
