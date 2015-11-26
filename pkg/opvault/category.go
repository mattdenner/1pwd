package opvault

type Category string

const (
	LoginItem           = Category("001")
	CreditCardItem      = Category("002")
	SecureNoteItem      = Category("003")
	IdentityItem        = Category("004")
	PasswordItem        = Category("005")
	TombstoneItem       = Category("099")
	SoftwareLicenseItem = Category("100")
	BankAccountItem     = Category("101")
	DatabaseItem        = Category("102")
	DriverLicenseItem   = Category("103")
	OutdoorLicenseItem  = Category("104")
	MembershipItem      = Category("105")
	PassportItem        = Category("106")
	RewardsItem         = Category("107")
	SSNItem             = Category("108")
	RouterItem          = Category("109")
	ServerItem          = Category("110")
	EmailItem           = Category("111")
)

func (c Category) String() string {
	switch c {
	case LoginItem:
		return "Login"
	case CreditCardItem:
		return "Credit Card"
	case SecureNoteItem:
		return "Secure Note"
	case IdentityItem:
		return "Identity"
	case PasswordItem:
		return "Password"
	case TombstoneItem:
		return "Tombstone"
	case SoftwareLicenseItem:
		return "Software License"
	case BankAccountItem:
		return "Bank Account"
	case DatabaseItem:
		return "Database"
	case DriverLicenseItem:
		return "Driver License"
	case OutdoorLicenseItem:
		return "Outdoor License"
	case MembershipItem:
		return "Membership"
	case PassportItem:
		return "Passport"
	case RewardsItem:
		return "Rewards"
	case SSNItem:
		return "SSN"
	case RouterItem:
		return "Router"
	case ServerItem:
		return "Server"
	case EmailItem:
		return "Email"
	default:
		return "Unknown"
	}
}

func (c Category) TypeString() string {
	switch c {
	case LoginItem:
		return "login"
	case CreditCardItem:
		return "credit-card"
	case SecureNoteItem:
		return "secure-note"
	case IdentityItem:
		return "identity"
	case PasswordItem:
		return "password"
	case TombstoneItem:
		return "tombstone"
	case SoftwareLicenseItem:
		return "software-license"
	case BankAccountItem:
		return "bank account"
	case DatabaseItem:
		return "database"
	case DriverLicenseItem:
		return "driver-license"
	case OutdoorLicenseItem:
		return "outdoor-license"
	case MembershipItem:
		return "membership"
	case PassportItem:
		return "passport"
	case RewardsItem:
		return "rewards"
	case SSNItem:
		return "ssn"
	case RouterItem:
		return "router"
	case ServerItem:
		return "server"
	case EmailItem:
		return "email"
	default:
		return "unknown"
	}
}

func FromTypeString(str string) Category {
	switch str {
	case "login":
		return LoginItem
	case "credit-card":
		return CreditCardItem
	case "secure-note":
		return SecureNoteItem
	case "identity":
		return IdentityItem
	case "password":
		return PasswordItem
	case "tombstone":
		return TombstoneItem
	case "software-license":
		return SoftwareLicenseItem
	case "bank account":
		return BankAccountItem
	case "database":
		return DatabaseItem
	case "driver-license":
		return DriverLicenseItem
	case "outdoor-license":
		return OutdoorLicenseItem
	case "membership":
		return MembershipItem
	case "passport":
		return PassportItem
	case "rewards":
		return RewardsItem
	case "ssn":
		return SSNItem
	case "router":
		return RouterItem
	case "server":
		return ServerItem
	case "email":
		return EmailItem
	default:
		return Category("Unknown")
	}
}
