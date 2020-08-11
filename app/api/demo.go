package api

//  // Demo 接口
//  type Demo struct {
//  	service.GPA
//  	DemoService *service.Demo
//  }
//
//  // Register 注册路由
//  func (a *Demo) Register(r gin.IRouter) {
//  	r.GET("hello", a.Hello)
//  	r.GET("get", a.Get)
//  	r.POST("set", a.Set)
//  	r.GET("get1", a.Get1)
//  }
//
//  // @Tags demo
//  // @Summary Hello
//  // @Description Hello world
//  // @Accept  json
//  // @Produce  json
//  // @Success 200 {string} string "ok"
//  // @Router /demo/hello [get]
//  func (a *Demo) hello(c *gin.Context) {
//  	c.JSON(http.StatusOK, gin.H{
//  		"message": "hello, world",
//  	})
//  }
//
//  // @Tags demo
//  // @Summary Set
//  // @Description Set
//  // @Accept  json
//  // @Produce  json
//  // @Param item body schema.DemoSet false "Demo Info"
//  // @Success 200 {object} helper.ErrorInfo
//  // @@Failure 500 {object} helper.ErrorInfo
//  // @Router /demo/set [post]
//  func (a *Demo) set(c *gin.Context) {
//  	item := &schema.DemoSet{}
//  	err := helper.ParseJSON(c, item)
//  	if err != nil {
//  		return
//  	}
//  	res, err := a.Entc.Demo.Create().
//  		SetCode(item.Code).
//  		SetName(item.Name).
//  		SetMemo(item.Memo).
//  		SetStatus(1).
//  		Save(c)
//
//  	if err != nil {
//  		panic(err)
//  	}
//  	helper.ResSuccess(c, res)
//  }
//
//  // @Tags demo
//  // @Summary Get
//  // @Description Get
//  // @Accept  json
//  // @Produce  json
//  // @Param id query string true "Demo id"
//  // @Success 200 {object} helper.ErrorInfo
//  // @@Failure 500 {object} helper.ErrorInfo
//  // @Router /demo/get [get]
//  func (a *Demo) get(c *gin.Context) {
//  	idstr := c.Query("id")
//  	if idstr == "" {
//  		helper.ResError(c, helper.Err406NotAcceptable)
//  		return
//  	}
//  	id, err := strconv.Atoi(idstr)
//  	if err != nil {
//  		helper.ResError(c, helper.Err406NotAcceptable)
//  		return
//  	}
//  	res, err := a.Entc.Demo.Get(c, id)
//  	if err != nil && !ent.IsNotFound(err) {
//  		panic(err)
//  	}
//
//  	helper.ResSuccess(c, res)
//  }
//
//  // @Tags demo
//  // @Summary Get
//  // @Description Get
//  // @Accept  json
//  // @Produce  json
//  // @Param id query string true "Demo id"
//  // @Success 200 {object} helper.ErrorInfo
//  // @@Failure 500 {object} helper.ErrorInfo
//  // @Router /demo/get1 [get]
//  func (a *Demo) get1(c *gin.Context) {
//  	idstr := c.Query("id")
//  	if idstr == "" {
//  		helper.ResError(c, helper.Err406NotAcceptable)
//  		return
//  	}
//  	id, err := strconv.Atoi(idstr)
//  	if err != nil {
//  		helper.ResError(c, helper.Err406NotAcceptable)
//  		return
//  	}
//  	res := &sqlxm.Demo{}
//  	err = a.Sqlx.Get(res, "SELECT id, code, name, memo FROM demo WHERE id=$1", id)
//  	if err != nil {
//  		if !sqlxc.IsNotFound(err) {
//  			panic(err)
//  		}
//  		res = nil // 没有找到数据
//  	}
//
//  	helper.ResSuccess(c, res)
//  }
