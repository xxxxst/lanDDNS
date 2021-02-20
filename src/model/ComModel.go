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
	DomainMatch			string	`ini:"domainMatch"`			;
	DefaultDnsServer1	string	`ini:"defaultDnsServer1"`	;
	DefaultDnsServer2	string	`ini:"defaultDnsServer2"`	;

	UseStaticHost		bool	`ini:"useStaticHost"`		;
	// UseDynamicHost		bool	`ini:"useDynamicHost"`		;
	UseDynamicHost		bool;
	UseMacHost			bool	`ini:"useMacHost"`			;
	MacIp				string	`ini:"macIp"`				;
}

type ComConfigComClient struct {
	// ManageIp		string	`ini:"manageIp"`		;
	// ManagePort		string	`ini:"managePort"`		;
	// EncryptTransfer	bool	`ini:"encryptTransfer"`	;
	ManageIp		string;
	ManagePort		string;
	EncryptTransfer	bool;
}

type ComConfigClient struct {
	// Domain			string	`ini:"domain"`			;
	Domain			string;
}

type ComConfig struct {
	Server			ComConfigServer		`ini:"server"`;
	// ComClient		ComConfigComClient	`ini:"comClient"`;
	ComClient		ComConfigComClient;
	// Client			ComConfigClient		`ini:"client"`;
	Client			ComConfigClient;
}

func CreateComConfig() (*ComConfig) {
	return &ComConfig{
		Server: ComConfigServer {
			DnsIp: "",
			DnsPort: "53",
			DomainMatch: "*.lan",
			DefaultDnsServer1: "114.114.114.114",
			DefaultDnsServer2: "8.8.8.8",

			UseStaticHost: true,
			UseDynamicHost: true,
			UseMacHost: true,
			MacIp: "",
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
