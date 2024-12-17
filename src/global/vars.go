package global

type DDNSUConfig struct {
	Ddnsu    Ddnsu    `toml:"ddnsu"`
	Services Services `toml:"services"`
}

type Ddnsu struct {
	Version     string   `toml:"version"`
	Use         string   `toml:"use"`
	IpProviders []string `toml:"ipProviders"`
	Rate        int      `toml:"rate"`
	Record      []Record `toml:"record"`
	Domain      string   `toml:"domain"`
}

type Services struct {
	Cloudflare Cloudflare `toml:"cloudflare"`
	Vercel     Vercel     `toml:"vercel"`
}

type Cloudflare struct {
	Token string `toml:"token"`
}

type Vercel struct {
	Token string `toml:"token"`
}

type Record struct {
	Rtype     string `toml:"rtype"`
	Comment   string `toml:"comment"`
	Ttl       int    `toml:"ttl"`
	Subdomain string `toml:"subdomain"`
}

var Token string = ""
var Configuration DDNSUConfig
var ConfigurationPath string
var CloudflareZoneId string
var LastIpAddress string

const RecordManagedPrefix = "[d]-"
