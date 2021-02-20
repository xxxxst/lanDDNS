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
	UseDynamicHost		bool	`ini:"useDynamicHost"`		;
	UseMacHost			bool	`ini:"useMacHost"`			;
	MacIp				string	`ini:"macIp"`				;
}

type ComConfigComClient struct {
	ManageIp		string	`ini:"manageIp"`		;
	ManagePort		string	`ini:"managePort"`		;
	EncryptTransfer	bool	`ini:"encryptTransfer"`	;
}

type ComConfigClient struct {
	Domain			string	`ini:"domain"`			;
}

type ComConfig struct {
	Server			ComConfigServer		`ini:"server"`;
	ComClient		ComConfigComClient	`ini:"comClient"`;
	Client			ComConfigClient		`ini:"client"`;
}

func CreateComConfig() (*ComConfig) {
	return &ComConfig{
		Server: ComConfigServer {
			DnsIp: "",
			DnsPort: "53",
			DomainMatch: "*.lan",
			DefaultDnsServer1: "8.8.8.8",
			DefaultDnsServer2: "114.114.114.114",

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
