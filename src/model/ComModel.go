package model

type ComModel struct {
	Version			string;
	IsDebug			bool;
	
	RootDir			string;
	ConfigDir		string;

	ConfigMd		*ComConfig;

	// DnsIp		string;
	// DnsPort		string;

	// ManageIp		string;
	// ManagePort		string;
}

var GetComModel = (func() (func() (*ComModel)) {
	var ins *ComModel;

	return func() (*ComModel) {
		if(ins == nil) {
			ins = new(ComModel);
			// ins.DnsIp = "0.0.0.0";
			// ins.DnsPort = "53";
			// ins.ManageIp = "0.0.0.0";
			// ins.ManagePort = "53";
		}
		return ins;
	}
})();

type ComConfigServer struct {
	DnsIp				string	`ini:"dnsIp"`				;
	DnsPort				string	`ini:"dnsPort"`				;
	// DomainMatch			string	`ini:"domainMatch"`			;
	DomainMatch			string	`ini:"-"`			;
	DefaultDnsServer1	string	`ini:"defaultDnsServer1"`	;
	DefaultDnsServer2	string	`ini:"defaultDnsServer2"`	;

	UseStaticHost		bool	`ini:"useStaticHost"`		;
	// UseDynamicHost		bool	`ini:"useDynamicHost"`		;
	UseDynamicHost		bool	`ini:"-"`		;
	UseMacHost			bool	`ini:"useMacHost"`			;
	MacIp				string	`ini:"macIp" comment:"network segment to watch, multi split with ',', mask is always 255.255.255.0"`				;
	LogMac				bool	`ini:"logMac"`				;
}

type ComConfigComClient struct {
	// ManageIp		string	`ini:"manageIp"`		;
	// ManagePort		string	`ini:"managePort"`		;
	// EncryptTransfer	bool	`ini:"encryptTransfer"`	;
	ManageIp		string	`ini:"-"`	;
	ManagePort		string	`ini:"-"`	;
	EncryptTransfer	bool	`ini:"-"`	;
}

type ComConfigClient struct {
	// Domain			string	`ini:"domain"`			;
	Domain			string	`ini:"-"`			;
}

type ComConfig struct {
	Server			ComConfigServer		`ini:"server"`;
	// ComClient		ComConfigComClient	`ini:"comClient"`;
	ComClient		ComConfigComClient	`ini:"-"`;
	// Client			ComConfigClient		`ini:"client"`;
	Client			ComConfigClient		`ini:"-"`;
}

func CreateComConfig() (*ComConfig) {
	return &ComConfig{
		Server: ComConfigServer {
			DnsIp: "0.0.0.0",
			DnsPort: "53",
			DomainMatch: "*.lan",
			DefaultDnsServer1: "114.114.114.114",
			DefaultDnsServer2: "8.8.8.8",

			UseStaticHost: true,
			UseDynamicHost: false,
			UseMacHost: true,
			MacIp: "192.168.0.0",
			LogMac: false,
		},

		ComClient: ComConfigComClient {
			ManageIp: "",
			ManagePort: "53",
			EncryptTransfer: true,
		},

		Client: ComConfigClient {
			Domain: "",
		},
	};
}
