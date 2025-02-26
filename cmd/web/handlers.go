package main

// 定义所有的处理器

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"SnippetBox.mikudayo.net/internal/models"

	"github.com/julienschmidt/httprouter"
)

// 存储用户输入的消息
type snippetCreateForm struct {
	// 告诉解码器去html里找name为`...`的input标签
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// 将验证器注入要验证的数据中
	// (类似于继承直接使当前要检验的数据结构拥有验证器所有的字段与方法)
	models.Validator `form:"-"`
}

// 存储用户填写的个人信息
type userSignupForm struct {
	Name             string `form:"name"`
	Email            string `form:"email"`
	Password         string `form:"password"`
	models.Validator `form:"-"`
}

// 存储用户的登录信息
type userLoginForm struct {
	Email            string `form:"email"`
	Password         string `form:"password"`
	models.Validator `form:"-"`
}

// 展示网站的主页面
func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	// 指定"/"的逻辑 防止预料之外的访问
	if r.URL.Path != "/" {
		// 直接使用not found方法写入信息
		app.notFound(w)
		return
	}

	// 每一个请求都是一个独立的goroutine 这里的panic不会导致server崩溃
	// panic("something go wrong...")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
	}
	// 先初始化默认数据再初始化查询得到的数据
	data := app.newTemplateData(r)
	data.Snippets = snippets
	// 使用render()进行home.tmpl.html渲染
	app.render(w, http.StatusOK, "home.tmpl.html", data)
	// w.Write([]byte("mikudayoooo"))
}

// 展示一个具体的消息页面
func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	// 获取url ?后的id用于数据库查询
	// 抽取为string并用于转换判断是否为正值
	// http://localhost:3939/snippet/view?id=39
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// w.Write([]byte("Display a specific miku..."))
	// 使用新的方法获取url中的值
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 0 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		// 判断是否是ErrNoRecord 这里要通过包名调用自己定义的错误model.ErrNoRecord
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}
	// 初始化用于渲染网页的信息
	// 取出临时存在于ctx中的数据并删除(一次性使用) 在这里如果信息不存在就会返回空的字符串
	// flash := app.sessionManager.PopString(r.Context(), "flash")
	data := app.newTemplateData(r)
	data.Snippet = snippet
	// 将创建成功的消息导入数据用于网页渲染
	// data.Flash = flash
	app.render(w, http.StatusOK, "view.tmpl.html", data)
	// 将搜索到的内容直接输出到响应体
	// fmt.Fprintf(w, "Display a specific miku %v...", snippet)
}

// 展示创建消息的页面
func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// 指定request的方法必须是POST
	// if r.Method != http.MethodPost {
	// 	// 使用定义在net/http中的常量来代替直接的字符串

	// 	// 添加新的header展示更详细的信息
	// 	w.Header().Set("Allow", http.MethodPost)

	// 	// 写入状态码405 表示request方法不被允许
	// 	// w.WriteHeader(http.StatusMethodNotAllowed)
	// 	// w.Write([]byte("method not allowed..."))

	// 	// 一般不会直接调用writeheader与write而是通过别的函数间接调用
	// 	// http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	// 使用三方库可以在初始化处理器的时候直接制定请求方法
	// 后续用于返回一个填写信息的网页

	// 创建数据结构体用于后续写入 因为这里需要数据中的时间(只初始化了时间)所以才需要再初始化这个结构体
	// 这也保证了代码的结构统一
	data := app.newTemplateData(r)
	form := snippetCreateForm{
		// 处理错误内容返回原网页重新填充的逻辑需要用到结构体存储信息
		// 在这里初始化初次进入页面看到的内容 如果没有设置这个结构体会因为尝试访问不存在的信息报错
		Expires: 365,
	}
	data.Form = form
	app.render(w, http.StatusOK, "create.tmpl.html", data)
	// w.Write([]byte("Createa a new miku..."))

}

// 获取用户在网页中填写的信息
func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// map必须字段也要进行初始化 否则会直接panic(assignment to entry in nil map)
	var form snippetCreateForm
	// 使用自定的helper公式化解码数据
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Get方法在查找不到数据的情况下是会返回空字符串的
	// 验证从用户端得到的信息是否正确

	// 创建一个map用于存储各种类型的错误
	// 确定title不是空值并且长度小于100 如果失败就直接把错误信息加入map
	form.CheckField(form.NotBlank(form.Title), "title", "标题不能为空...")
	form.CheckField(form.MaxChars(form.Title, 100), "title", "标题长度不能超过100个字符...")
	form.CheckField(form.NotBlank(form.Content), "content", "内容不能为空...")
	form.CheckField(models.PermittedValue(form.Expires, 3, 7, 365), "expires", "时间必须为1,7,365...")
	// 检测是否有字段出现错误
	if !form.Valid() {
		// 如果有字段出现错误就以原先的输入信息重新渲染网页
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// curl -iL -X POST http://localhost:3939/snippet/create
	// 创建成功后为当前用户的会话添加共享信息(如果key存在则会将原先的信息覆盖掉)
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	// 创建成功后将用户重定向到最新创建的snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

// 展示用户的注册页面
func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// 初始话默认结构体为空值 防止因未传入默认结构体导致网页初始化错误
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

// 将用户填写的注册信息发送到后端
func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	// 定义用于存储输入信息的结构体并尝试向其解码
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// 解码成功检查数据的正确性
	form.CheckField(form.NotBlank(form.Name), "name", "姓名不能为空...")
	form.CheckField(form.NotBlank(form.Email), "email", "邮箱不能为空...")
	form.CheckField(form.Matches(form.Email, models.EmailRX), "email", "输入的邮箱格式错误...")
	form.CheckField(form.MinChars(form.Password, 8), "password", "密码长度必须大于8...")
	// 如果填入的字段出现错误就将字段返回给网页重新渲染
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}
	// 没有出现错误
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		// 判断错误类型
		if errors.Is(err, models.ErrDuplicateEmail) {
			// 在这里返回的已经是将原始错误包装过形成的自定义错误
			form.AddFieldError("email", "输入的邮箱已经被使用...")

			// 为用户重新渲染页面
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		// 终止请求
		return
	}
	// 创建成功了使用session创建flash信息进行提示
	app.sessionManager.Put(r.Context(), "flash", "注册成功！请登入...")

	// 对网页进行重定向
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// 展示用户的登录界面
func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	// 初始化参数用于网页渲染
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
	// fmt.Fprint(w, "Display a html form for logging in a user...")
}

// 将用户填写的登录信息发送到后端
func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	// 尝试解码用户填写的数据
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// 先对输入的数据进行简单的有效性验证
	form.CheckField(form.NotBlank(form.Email), "email", "邮箱不能为空值...")
	form.CheckField(form.Matches(form.Email, models.EmailRX), "email", "邮箱格式不正确...")
	form.CheckField(form.NotBlank(form.Password), "password", "密码不能为空值...")
	// 如果有错误返回填写的信息重新渲染网页
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}
	// 填写信息的格式都正确进行正式的有效性检测
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		// 判断错误是否是无效数据错误
		if errors.Is(err, models.ErrInvalidCredentials) {
			// 将错误信息添加到NonFieldErrors
			form.AddNonFieldError("邮箱或密码错误...")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// 在登陆成功后或者权限等级发生变化后更新session ID
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// 验证通过将当前用户的id加入session表示已登入
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	// 重定向到创建消息页面表示当前已登录
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
	// fmt.Fprint(w, "Authenticate and login the user...")
}

// 将用户需要退出的信息发送到后端
func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// 更新会话ID
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	// 删除当前登入的AuthenticateUserID
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	// 添加一个flash消息提示用户成功退出
	app.sessionManager.Put(r.Context(), "flash", "已成功退出...")
	// 导航回到主页面
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
