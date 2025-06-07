package utils

import (
	"net/url"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID      `bson:"_id,omitempty"`
	DiscordID string                  `bson:"discord_id"`
	Username  string                  `bson:"username"`
	Email     string                  `bson:"email"`
	Guild     *GuildInfo              `bson:"guild"`
	Channels  map[string]*ChannelInfo `bson:"channels"`
	LeavedAt  time.Time               `bson:"leaved_at"`
}

type GuildInfo struct {
	GuildID  string    `bson:"guild_id"`
	JoinedAt time.Time `bson:"joined_at"`
}

type ChannelInfo struct {
	Name  string            `bson:"name"`
	Links map[string]string `bson:"links,omitempty"`
}

type AuthCookie struct {
	Accesstoken struct {
		String  string
		Expires time.Duration
	}
	Refreshtoken struct {
		String  string
		Expires time.Duration
	}
}

type RequestBuilder struct {
	URL    *url.URL
	Method string
	Cookie *AuthCookie
	Proxy  *url.URL
}

type DiscordUserData struct {
	GuildID     string
	Member      *discordgo.Member
	ChannelID   string
	ChannelName string
	Link        string
	LinkName    string
}

type SessionData struct {
	Channels map[string]*ChannelInfo `bson:"channels"`
}

type ItemsResp struct {
	Items []CatalogItem `json:"items"`
}

type CatalogItem struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Price struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"price"`
	IsVisible  bool        `json:"is_visible"`
	Discount   interface{} `json:"discount"`
	BrandTitle string      `json:"brand_title"`
	Path       string      `json:"path"`
	User       struct {
		ID         int    `json:"id"`
		Login      string `json:"login"`
		ProfileURL string `json:"profile_url"`
		Photo      struct {
			ID                  int         `json:"id"`
			Width               int         `json:"width"`
			Height              int         `json:"height"`
			TempUUID            interface{} `json:"temp_uuid"`
			URL                 string      `json:"url"`
			DominantColor       string      `json:"dominant_color"`
			DominantColorOpaque string      `json:"dominant_color_opaque"`
			Thumbnails          []struct {
				Type         string      `json:"type"`
				URL          string      `json:"url"`
				Width        int         `json:"width"`
				Height       int         `json:"height"`
				OriginalSize interface{} `json:"original_size"`
			} `json:"thumbnails"`
			IsSuspicious   bool        `json:"is_suspicious"`
			Orientation    interface{} `json:"orientation"`
			HighResolution struct {
				ID          string      `json:"id"`
				Timestamp   int         `json:"timestamp"`
				Orientation interface{} `json:"orientation"`
			} `json:"high_resolution"`
			FullSizeURL string `json:"full_size_url"`
			IsHidden    bool   `json:"is_hidden"`
			Extra       struct {
			} `json:"extra"`
		} `json:"photo"`
		Business bool `json:"business"`
	} `json:"user"`
	Conversion interface{} `json:"conversion"`
	URL        string      `json:"url"`
	Promoted   bool        `json:"promoted"`
	Photo      struct {
		ID                  int64  `json:"id"`
		ImageNo             int    `json:"image_no"`
		Width               int    `json:"width"`
		Height              int    `json:"height"`
		DominantColor       string `json:"dominant_color"`
		DominantColorOpaque string `json:"dominant_color_opaque"`
		URL                 string `json:"url"`
		IsMain              bool   `json:"is_main"`
		Thumbnails          []struct {
			Type         string      `json:"type"`
			URL          string      `json:"url"`
			Width        int         `json:"width"`
			Height       int         `json:"height"`
			OriginalSize interface{} `json:"original_size"`
		} `json:"thumbnails"`
		HighResolution struct {
			ID          string `json:"id"`
			Timestamp   int    `json:"timestamp"`
			Orientation int    `json:"orientation"`
		} `json:"high_resolution"`
		IsSuspicious bool   `json:"is_suspicious"`
		FullSizeURL  string `json:"full_size_url"`
		IsHidden     bool   `json:"is_hidden"`
		Extra        struct {
		} `json:"extra"`
	} `json:"photo"`
	FavouriteCount int         `json:"favourite_count"`
	IsFavourite    bool        `json:"is_favourite"`
	Badge          interface{} `json:"badge"`
	ServiceFee     struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"service_fee"`
	TotalItemPrice struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"total_item_price"`
	ViewCount     int           `json:"view_count"`
	SizeTitle     string        `json:"size_title"`
	ContentSource string        `json:"content_source"`
	Status        string        `json:"status"`
	IconBadges    []interface{} `json:"icon_badges"`
	ItemBox       struct {
		FirstLine          string `json:"first_line"`
		SecondLine         string `json:"second_line"`
		AccessibilityLabel string `json:"accessibility_label"`
	} `json:"item_box"`
	SearchTrackingParams struct {
		Score          int           `json:"score"`
		MatchedQueries []interface{} `json:"matched_queries"`
	} `json:"search_tracking_params"`
}

// Deprecated
type Item struct {
	ID                           int64         `json:"id"`
	Title                        string        `json:"title"`
	BrandID                      int           `json:"brand_id"`
	SizeID                       interface{}   `json:"size_id"`
	StatusID                     int           `json:"status_id"`
	UserID                       int           `json:"user_id"`
	CountryID                    int           `json:"country_id"`
	CatalogID                    int           `json:"catalog_id"`
	Color1ID                     int           `json:"color1_id"`
	Color2ID                     int           `json:"color2_id"`
	PackageSizeID                int           `json:"package_size_id"`
	IsUnisex                     int           `json:"is_unisex"`
	ModerationStatus             int           `json:"moderation_status"`
	IsHidden                     bool          `json:"is_hidden"`
	IsHiddenWithItemRestrictions bool          `json:"is_hidden_with_item_restrictions"`
	IsVisible                    bool          `json:"is_visible"`
	IsClosed                     bool          `json:"is_closed"`
	FavouriteCount               int           `json:"favourite_count"`
	ActiveBidCount               int           `json:"active_bid_count"`
	Description                  string        `json:"description"`
	PackageSizeStandard          bool          `json:"package_size_standard"`
	ItemClosingAction            interface{}   `json:"item_closing_action"`
	RelatedCatalogIds            []interface{} `json:"related_catalog_ids"`
	RelatedCatalogsEnabled       bool          `json:"related_catalogs_enabled"`
	Size                         string        `json:"size"`
	Brand                        string        `json:"brand"`
	Composition                  string        `json:"composition"`
	ExtraConditions              string        `json:"extra_conditions"`
	DisposalConditions           int           `json:"disposal_conditions"`
	IsForSell                    bool          `json:"is_for_sell"`
	IsHandicraft                 bool          `json:"is_handicraft"`
	IsProcessing                 bool          `json:"is_processing"`
	IsDraft                      bool          `json:"is_draft"`
	IsReserved                   bool          `json:"is_reserved"`
	Label                        string        `json:"label"`
	OriginalPriceNumeric         string        `json:"original_price_numeric"`
	Currency                     string        `json:"currency"`
	PriceNumeric                 string        `json:"price_numeric"`
	LastPushUpAt                 time.Time     `json:"last_push_up_at"`
	CreatedAtTs                  time.Time     `json:"created_at_ts"`
	UpdatedAtTs                  time.Time     `json:"updated_at_ts"`
	UserUpdatedAtTs              time.Time     `json:"user_updated_at_ts"`
	IsDelayedPublication         bool          `json:"is_delayed_publication"`
	Photos                       []struct {
		ID                  int64  `json:"id"`
		ImageNo             int    `json:"image_no"`
		Width               int    `json:"width"`
		Height              int    `json:"height"`
		DominantColor       string `json:"dominant_color"`
		DominantColorOpaque string `json:"dominant_color_opaque"`
		URL                 string `json:"url"`
		IsMain              bool   `json:"is_main"`
		Thumbnails          []struct {
			Type         string      `json:"type"`
			URL          string      `json:"url"`
			Width        int         `json:"width"`
			Height       int         `json:"height"`
			OriginalSize interface{} `json:"original_size"`
		} `json:"thumbnails"`
		HighResolution struct {
			ID          string `json:"id"`
			Timestamp   int    `json:"timestamp"`
			Orientation int    `json:"orientation"`
		} `json:"high_resolution"`
		IsSuspicious bool   `json:"is_suspicious"`
		FullSizeURL  string `json:"full_size_url"`
		IsHidden     bool   `json:"is_hidden"`
		Extra        struct {
		} `json:"extra"`
	} `json:"photos"`
	CanBeSold               bool        `json:"can_be_sold"`
	CanFeedback             bool        `json:"can_feedback"`
	ItemReservationID       interface{} `json:"item_reservation_id"`
	PromotedUntil           interface{} `json:"promoted_until"`
	PromotedInternationally interface{} `json:"promoted_internationally"`
	DiscountPriceNumeric    interface{} `json:"discount_price_numeric"`
	Author                  interface{} `json:"author"`
	BookTitle               interface{} `json:"book_title"`
	Isbn                    interface{} `json:"isbn"`
	MeasurementWidth        interface{} `json:"measurement_width"`
	MeasurementLength       interface{} `json:"measurement_length"`
	MeasurementUnit         interface{} `json:"measurement_unit"`
	Manufacturer            interface{} `json:"manufacturer"`
	ManufacturerLabelling   interface{} `json:"manufacturer_labelling"`
	TransactionPermitted    bool        `json:"transaction_permitted"`
	VideoGameRatingID       interface{} `json:"video_game_rating_id"`
	ItemAttributes          []struct {
		Code string `json:"code"`
		Ids  []int  `json:"ids"`
	} `json:"item_attributes"`
	HaovItem     bool `json:"haov_item?"`
	CollectionID int  `json:"collection_id"`
	ModelID      int  `json:"model_id"`
	Price        struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"price"`
	DiscountPrice interface{} `json:"discount_price"`
	User          struct {
		ID                 int         `json:"id"`
		AnonID             string      `json:"anon_id"`
		Login              string      `json:"login"`
		Birthday           interface{} `json:"birthday"`
		CountryID          int         `json:"country_id"`
		CityID             int         `json:"city_id"`
		UpdatedOn          int         `json:"updated_on"`
		RealName           interface{} `json:"real_name"`
		Email              interface{} `json:"email"`
		CountryCode        string      `json:"country_code"`
		FeedbackCount      int         `json:"feedback_count"`
		FeedbackReputation float64     `json:"feedback_reputation"`
		Moderator          bool        `json:"moderator"`
		BusinessAccountID  interface{} `json:"business_account_id"`
		CanBundle          bool        `json:"can_bundle"`
		BusinessAccount    interface{} `json:"business_account"`
		Business           bool        `json:"business"`
		Photo              struct {
			ID                  int         `json:"id"`
			Width               int         `json:"width"`
			Height              int         `json:"height"`
			TempUUID            interface{} `json:"temp_uuid"`
			URL                 string      `json:"url"`
			DominantColor       string      `json:"dominant_color"`
			DominantColorOpaque string      `json:"dominant_color_opaque"`
			Thumbnails          []struct {
				Type         string      `json:"type"`
				URL          string      `json:"url"`
				Width        int         `json:"width"`
				Height       int         `json:"height"`
				OriginalSize interface{} `json:"original_size"`
			} `json:"thumbnails"`
			IsSuspicious   bool        `json:"is_suspicious"`
			Orientation    interface{} `json:"orientation"`
			HighResolution struct {
				ID          string      `json:"id"`
				Timestamp   int         `json:"timestamp"`
				Orientation interface{} `json:"orientation"`
			} `json:"high_resolution"`
			FullSizeURL string `json:"full_size_url"`
			IsHidden    bool   `json:"is_hidden"`
			Extra       struct {
			} `json:"extra"`
		} `json:"photo"`
		AcceptedPayInMethods []struct {
			ID                   int    `json:"id"`
			Code                 string `json:"code"`
			RequiresCreditCard   bool   `json:"requires_credit_card"`
			EventTrackingCode    string `json:"event_tracking_code"`
			Icon                 string `json:"icon"`
			Enabled              bool   `json:"enabled"`
			TranslatedName       string `json:"translated_name"`
			Note                 string `json:"note"`
			MethodChangePossible bool   `json:"method_change_possible"`
		} `json:"accepted_pay_in_methods"`
		CanViewProfile               bool          `json:"can_view_profile"`
		BundleDiscount               interface{}   `json:"bundle_discount"`
		LastLogedOnTs                time.Time     `json:"last_loged_on_ts"`
		LastLogedOn                  string        `json:"last_loged_on"`
		ItemCount                    int           `json:"item_count"`
		TotalItemsCount              int           `json:"total_items_count"`
		FollowersCount               int           `json:"followers_count"`
		FollowingCount               int           `json:"following_count"`
		FollowingBrandsCount         int           `json:"following_brands_count"`
		AccountStatus                int           `json:"account_status"`
		PositiveFeedbackCount        int           `json:"positive_feedback_count"`
		NeutralFeedbackCount         int           `json:"neutral_feedback_count"`
		NegativeFeedbackCount        int           `json:"negative_feedback_count"`
		IsOnHoliday                  bool          `json:"is_on_holiday"`
		ExposeLocation               bool          `json:"expose_location"`
		City                         string        `json:"city"`
		IsPublishPhotosAgreed        bool          `json:"is_publish_photos_agreed"`
		ThirdPartyTracking           bool          `json:"third_party_tracking"`
		Locale                       string        `json:"locale"`
		ProfileURL                   string        `json:"profile_url"`
		ShareProfileURL              string        `json:"share_profile_url"`
		IsOnline                     bool          `json:"is_online"`
		Fundraiser                   interface{}   `json:"fundraiser"`
		Localization                 string        `json:"localization"`
		IsBpfPriceProminenceApplied  bool          `json:"is_bpf_price_prominence_applied"`
		MsgTemplateCount             int           `json:"msg_template_count"`
		IsAccountBanned              bool          `json:"is_account_banned"`
		AccountBanDate               interface{}   `json:"account_ban_date"`
		IsAccountBanPermanent        bool          `json:"is_account_ban_permanent"`
		SellerBadges                 []interface{} `json:"seller_badges"`
		IsFavourite                  bool          `json:"is_favourite"`
		IsHated                      bool          `json:"is_hated"`
		HatesYou                     bool          `json:"hates_you"`
		ContactsPermission           interface{}   `json:"contacts_permission"`
		Contacts                     interface{}   `json:"contacts"`
		Path                         string        `json:"path"`
		IsCatalogModerator           bool          `json:"is_catalog_moderator"`
		IsCatalogRoleMarketingPhotos bool          `json:"is_catalog_role_marketing_photos"`
		HideFeedback                 bool          `json:"hide_feedback"`
		AllowDirectMessaging         bool          `json:"allow_direct_messaging"`
		Verification                 struct {
			Email struct {
				Valid     bool `json:"valid"`
				Available bool `json:"available"`
			} `json:"email"`
			Facebook struct {
				Valid      bool        `json:"valid"`
				VerifiedAt interface{} `json:"verified_at"`
				Available  bool        `json:"available"`
			} `json:"facebook"`
			Google struct {
				Valid      bool        `json:"valid"`
				VerifiedAt interface{} `json:"verified_at"`
				Available  bool        `json:"available"`
			} `json:"google"`
		} `json:"verification"`
		AvgResponseTime   interface{} `json:"avg_response_time"`
		About             string      `json:"about"`
		FacebookUserID    interface{} `json:"facebook_user_id"`
		GivenItemCount    int         `json:"given_item_count"`
		TakenItemCount    int         `json:"taken_item_count"`
		CountryTitleLocal string      `json:"country_title_local"`
		CountryIsoCode    string      `json:"country_iso_code"`
		CountryTitle      string      `json:"country_title"`
		DefaultAddress    interface{} `json:"default_address"`
	} `json:"user"`
	OfflineVerification    bool `json:"offline_verification"`
	OfflineVerificationFee struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"offline_verification_fee"`
	ServiceFee       string `json:"service_fee"`
	TotalItemPrice   string `json:"total_item_price"`
	CanEdit          bool   `json:"can_edit"`
	CanDelete        bool   `json:"can_delete"`
	CanReserve       bool   `json:"can_reserve"`
	CanMarkAsSold    bool   `json:"can_mark_as_sold"`
	CanTransfer      bool   `json:"can_transfer"`
	InstantBuy       bool   `json:"instant_buy"`
	CanClose         bool   `json:"can_close"`
	CanBuy           bool   `json:"can_buy"`
	CanBundle        bool   `json:"can_bundle"`
	CanAskSeller     bool   `json:"can_ask_seller"`
	CanFavourite     bool   `json:"can_favourite"`
	UserLogin        string `json:"user_login"`
	CityID           int    `json:"city_id"`
	City             string `json:"city"`
	Country          string `json:"country"`
	Promoted         bool   `json:"promoted"`
	IsMobile         bool   `json:"is_mobile"`
	BumpBadgeVisible bool   `json:"bump_badge_visible"`
	BrandDto         struct {
		ID                        int    `json:"id"`
		Title                     string `json:"title"`
		Slug                      string `json:"slug"`
		FavouriteCount            int    `json:"favourite_count"`
		PrettyFavouriteCount      string `json:"pretty_favourite_count"`
		ItemCount                 int    `json:"item_count"`
		PrettyItemCount           string `json:"pretty_item_count"`
		IsVisibleInListings       bool   `json:"is_visible_in_listings"`
		RequiresAuthenticityCheck bool   `json:"requires_authenticity_check"`
		IsLuxury                  bool   `json:"is_luxury"`
		IsHvf                     bool   `json:"is_hvf"`
		Path                      string `json:"path"`
		URL                       string `json:"url"`
		IsFavourite               bool   `json:"is_favourite"`
	} `json:"brand_dto"`
	CatalogBranchTitle   string `json:"catalog_branch_title"`
	Path                 string `json:"path"`
	URL                  string `json:"url"`
	AcceptedPayInMethods []struct {
		ID                   int    `json:"id"`
		Code                 string `json:"code"`
		RequiresCreditCard   bool   `json:"requires_credit_card"`
		EventTrackingCode    string `json:"event_tracking_code"`
		Icon                 string `json:"icon"`
		Enabled              bool   `json:"enabled"`
		TranslatedName       string `json:"translated_name"`
		Note                 string `json:"note"`
		MethodChangePossible bool   `json:"method_change_possible"`
	} `json:"accepted_pay_in_methods"`
	CreatedAt             string `json:"created_at"`
	Color1                string `json:"color1"`
	Color2                string `json:"color2"`
	SizeTitle             string `json:"size_title"`
	DescriptionAttributes []struct {
		Code  string      `json:"code"`
		Title string      `json:"title"`
		Value string      `json:"value"`
		FaqID interface{} `json:"faq_id"`
	} `json:"description_attributes"`
	VideoGameRating     interface{}   `json:"video_game_rating"`
	Status              string        `json:"status"`
	IsFavourite         bool          `json:"is_favourite"`
	ViewCount           int           `json:"view_count"`
	Performance         interface{}   `json:"performance"`
	StatsVisible        bool          `json:"stats_visible"`
	CanPushUp           bool          `json:"can_push_up"`
	ItemAlert           interface{}   `json:"item_alert"`
	ItemAlertType       interface{}   `json:"item_alert_type"`
	Badge               interface{}   `json:"badge"`
	SizeGuideFaqEntryID interface{}   `json:"size_guide_faq_entry_id"`
	Localization        string        `json:"localization"`
	IconBadges          []interface{} `json:"icon_badges"`
	ItemBox             struct {
		FirstLine          string `json:"first_line"`
		SecondLine         string `json:"second_line"`
		AccessibilityLabel string `json:"accessibility_label"`
	} `json:"item_box"`
}
