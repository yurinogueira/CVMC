package auth

type Permission string

const (
	PermissionEditCar           Permission = "edit_car"
	PermissionDeleteCar         Permission = "delete_car"
	PermissionShareCar          Permission = "share_car"
	PermissionRemoveSharing     Permission = "remove_sharing"
	PermissionCreateMaintenance Permission = "create_maintenance"
	PermissionEditMaintenance   Permission = "edit_maintenance"
	PermissionDeleteMaintenance Permission = "delete_maintenance"
	PermissionViewCar           Permission = "view_car"
	PermissionViewHistory       Permission = "view_history"
	PermissionViewAttachments   Permission = "view_attachments"
)
