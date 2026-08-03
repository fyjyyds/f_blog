package service

// Container 持有所有 service 实例，供 handler 使用
type Container struct {
	User         *UserService
	Article      *ArticleService
	Comment      *CommentService
	Like         *LikeService
	Follow       *FollowService
	Notification *NotificationService
	Category     *CategoryService
	Tag          *TagService
	Banner       *BannerService
	Setting      *SettingService
}

var App *Container

func InitContainer(c *Container) {
	App = c
}
