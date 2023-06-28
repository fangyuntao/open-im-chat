package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/gin-gonic/gin"
)

func NewChatRoute(router gin.IRouter, zk discoveryregistry.SvcDiscoveryRegistry) {
	mw := NewMW(zk)
	chat := NewChat(zk)
	account := router.Group("/account")
	account.POST("/code/send", chat.SendVerifyCode)                      //发送验证码
	account.POST("/code/verify", chat.VerifyCode)                        //校验验证码
	account.POST("/register", chat.RegisterUser)                         //注册
	account.POST("/login", chat.Login)                                   //登录
	account.POST("/password/reset", chat.ResetPassword)                  //忘记密码
	account.POST("/password/change", mw.CheckToken, chat.ChangePassword) //修改密码

	user := router.Group("/user", mw.CheckToken)
	user.POST("/update", chat.UpdateUserInfo)              //编辑个人资料
	user.POST("/find/public", chat.FindUserPublicInfo)     //获取用户公开信息
	user.POST("/find/full", chat.FindUserFullInfo)         //获取用户所有信息
	user.POST("/search/full", chat.SearchUserFullInfo)     //搜索用户公开信息
	user.POST("/search/public", chat.SearchUserPublicInfo) //搜索用户所有信息

	router.Group("/applet").POST("/find", mw.CheckToken, chat.FindApplet) //小程序列表

	router.Group("/client_config").POST("/get", chat.GetClientConfig) //获取客户端初始化配置

	router.Group("/callback").POST("/open_im", chat.OpenIMCallback) //回调
}

func NewAdminRoute(router gin.IRouter, zk discoveryregistry.SvcDiscoveryRegistry) {
	mw := NewMW(zk)
	admin := NewAdmin(zk)

	adminRouterGroup := router.Group("/account")
	adminRouterGroup.POST("/login", admin.AdminLogin)                      //登录
	adminRouterGroup.POST("/update", mw.CheckAdmin, admin.AdminUpdateInfo) //修改信息
	adminRouterGroup.POST("/info", mw.CheckAdmin, admin.AdminInfo)         //获取信息

	defaultRouter := router.Group("/default")
	defaultUserRouter := defaultRouter.Group("/user", mw.CheckAdmin)
	defaultUserRouter.POST("/add", admin.AddDefaultFriend)       //添加注册时默认好友
	defaultUserRouter.POST("/del", admin.DelDefaultFriend)       //删除注册时默认好友
	defaultUserRouter.POST("/find", admin.FindDefaultFriend)     //默认好友列表
	defaultUserRouter.POST("/search", admin.SearchDefaultFriend) //搜索注册时默认好友列表
	defaultGroupRouter := defaultRouter.Group("/group")
	defaultGroupRouter.POST("/add", admin.AddDefaultGroup)       //添加注册时默认群
	defaultGroupRouter.POST("/del", admin.DelDefaultGroup)       //删除注册时默认群
	defaultGroupRouter.POST("/find", admin.FindDefaultGroup)     //获取注册时默认群列表
	defaultGroupRouter.POST("/search", admin.SearchDefaultGroup) //获取注册时默认群列表

	invitationCodeRouter := router.Group("/invitation_code", mw.CheckAdmin)
	invitationCodeRouter.POST("/add", admin.AddInvitationCode)       //添加邀请码
	invitationCodeRouter.POST("/gen", admin.GenInvitationCode)       //生成邀请码
	invitationCodeRouter.POST("/del", admin.DelInvitationCode)       //删除邀请码
	invitationCodeRouter.POST("/search", admin.SearchInvitationCode) //搜索邀请码

	forbiddenRouter := router.Group("/forbidden", mw.CheckAdmin)
	ipForbiddenRouter := forbiddenRouter.Group("/ip")
	ipForbiddenRouter.POST("/add", admin.AddIPForbidden)       //添加禁止注册登录IP
	ipForbiddenRouter.POST("/del", admin.DelIPForbidden)       //删除禁止注册登录IP
	ipForbiddenRouter.POST("/search", admin.SearchIPForbidden) //搜索禁止注册登录IP
	userForbiddenRouter := forbiddenRouter.Group("/user")
	userForbiddenRouter.POST("/add", admin.AddUserIPLimitLogin)       //添加限制用户在指定ip登录
	userForbiddenRouter.POST("/del", admin.DelUserIPLimitLogin)       //删除用户在指定IP登录
	userForbiddenRouter.POST("/search", admin.SearchUserIPLimitLogin) //搜索限制用户在指定ip登录

	appletRouterGroup := router.Group("/applet", mw.CheckAdmin)
	appletRouterGroup.POST("/add", admin.AddApplet)       //添加小程序
	appletRouterGroup.POST("/del", admin.DelApplet)       //删除小程序
	appletRouterGroup.POST("/update", admin.UpdateApplet) //修改小程序
	appletRouterGroup.POST("/search", admin.SearchApplet) //搜索小程序

	blockRouter := router.Group("/block", mw.CheckAdmin)
	blockRouter.POST("/add", admin.BlockUser)          //封号
	blockRouter.POST("/del", admin.UnblockUser)        //解封
	blockRouter.POST("/search", admin.SearchBlockUser) //搜索封号用户

	userRouter := router.Group("/user", mw.CheckAdmin)
	userRouter.POST("/password/reset", admin.ResetUserPassword) //重置用户密码

	initGroup := router.Group("/client_config", mw.CheckAdmin)
	initGroup.POST("/set", admin.SetClientConfig) //设置客户端初始化配置
	initGroup.POST("/get", admin.GetClientConfig) //获取客户端初始化配置
}
func NewOrganizationRoute(router gin.IRouter, zk discoveryregistry.SvcDiscoveryRegistry) {
	mw := NewMW(zk)
	org := NewOrg(zk)
	userGroup := router.Group("/user")
	{
		userGroup.POST("/reset_password", mw.CheckToken, org.UpdateUserPassword) // 修改密码
		userGroup.POST("/login", org.UserLogin)                                  // 登录
		//userGroup.POST("/get_token", org.GetUserToken)            // 使用管理员token或者secret获取用户token
		userGroup.POST("/update", mw.CheckUser, org.UpdateUserInfo)          // 修改用户信息
		userGroup.POST("/info", mw.CheckToken, org.GetUserInfo)              // 获取用户信息
		userGroup.POST("/delete", mw.CheckAdmin, org.DeleteOrganizationUser) // 删除用户

		userGroup.POST("/get_users_full_info", org.GetUserFullList)        // 获取用户信息
		userGroup.POST("/search_users_full_info", org.SearchUsersFullInfo) // 搜索用户信息

		userGroup.POST("/callback", org.Callback)
	}

	organizationGroup := router.Group("/organization")
	{
		//部门  增删改查
		organizationGroup.POST("/create_department", mw.CheckAdmin, org.CreateDepartment) // 创建部门
		organizationGroup.POST("/update_department", org.UpdateDepartment)                // 修改部门
		organizationGroup.POST("/delete_department", mw.CheckAdmin, org.DeleteDepartment) // 删除部门
		organizationGroup.POST("/get_department", mw.CheckToken, org.GetDepartment)       // 获取部门

		//用户 增删改查
		organizationGroup.POST("/create_organization_user", mw.CheckAdmin, org.CreateOrganizationUser) // 创建用户 在某个部门或公司中新增
		organizationGroup.POST("/update_organization_user", mw.CheckAdmin, org.UpdateOrganizationUser) // 修改用户信息
		organizationGroup.POST("/delete_organization_user", org.DeleteOrganizationUser)                // 删除用户

		//查询用户所在的部门信息以及个人资料
		organizationGroup.POST("/get_user_in_department", mw.CheckToken, org.GetUserInDepartment)       // 获取用户所在部门
		organizationGroup.POST("/create_department_member", mw.CheckAdmin, org.CreateDepartmentMember)  // 创建部门成员 在某个部门或公司中新增
		organizationGroup.POST("/update_user_in_department", mw.CheckAdmin, org.UpdateUserInDepartment) // 修改用户部门
		//删除
		//organizationGroup.POST("/get_department_member", organization.GetDepartmentMember)        // 获取部门成员
		organizationGroup.POST("/delete_user_in_department", mw.CheckAdmin, org.DeleteUserInDepartment) // 删除部门成员 批量

		organizationGroup.POST("/get_search_user", mw.CheckAdmin, org.GetSearchUserList) // 搜索列表 后端

		organizationGroup.POST("/set_organization", mw.CheckAdmin, org.SetOrganization)        // 设置公司信息
		organizationGroup.POST("/get_organization", mw.CheckToken, org.GetOrganization)        // 获取公司信息
		organizationGroup.POST("/move_user_department", mw.CheckAdmin, org.MoveUserDepartment) // 移动用户部门

		organizationGroup.POST("/get_sub_department", mw.CheckToken, org.GetSubDepartment) // 获取部门的人和同级部门

		organizationGroup.POST("/get_search_department_user", mw.CheckToken, org.GetSearchDepartmentUser) // 搜索部门和用户
		organizationGroup.POST("/get_search_department_user_without_token", org.GetSearchDepartmentUserWithoutToken)

		organizationGroup.POST("/get_organization_department", mw.CheckToken, org.GetOrganizationDepartment) // 获取组织部门

		organizationGroup.POST("/sort_department", mw.CheckAdmin, org.SortDepartmentList)
		organizationGroup.POST("/sort_organization_user", mw.CheckAdmin, org.SortOrganizationUserList)

		organizationGroup.POST("/create_new_organization_member", mw.CheckAdmin, org.CreateNewOrganizationMember) // 创建用户的同时为其添加部门

		organizationGroup.POST("/import", org.BatchImport)                 // 批量导入
		organizationGroup.GET("/import_template", org.BatchImportTemplate) // 批量导入模板
	}
}
