package integration

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"

	"k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCreateContainer(t *testing.T) {

	host := "https://192.168.99.100:8443"
	// host := "https://10.10.7.22:6443"
	// cacert := `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQVRBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwdGFXNXAKYTNWaVpVTkJNQjRYRFRFNE1ERXdPVEV5TlRFeE5Wb1hEVEk0TURFd056RXlOVEV4TlZvd0ZURVRNQkVHQTFVRQpBeE1LYldsdWFXdDFZbVZEUVRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTTYwCjd2SjNISTh3bG9XSWF2UkJneDRzWXROZGNLeFpiRlBWYk8rajRzemVVRlN1TllCUitXdmlHMzgyelQvKzVBNDgKT28wUHhib3g3djdRQXBua1RLTGNrT21xR05DN3RlOHNRVWtFemRSWGM3Y0x0Nmh4dHBrb2NraEtwK0ltM1BJQwpvQ2hzQmMwTWxWSmpiR25IRHJXZjFZa2JxaHQ1c3lBRTBEbVVGbXJ6MWdvdEgwN0d1Rkw2b0oyWjJCeXdmaERzCjh1cFB3djR6VzJLbXV4MjRwNUNnOEhVeTV4N2NmeFhQN2tsbkY5bVFRSDVqcXp6eDRBTnBWWURJTzA0QXZJK3cKQ1JXdGtlT3hpSmxXRkRFcjF0NURtcmp2YUlYcWdGN0ROYWlhMWtpcm94czFKUlFGM0xTK21va09iNmt2N1ZNVgpMa1dwZ0VMMG9DRnZybFMyT2k4Q0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUIwR0ExVWRKUVFXCk1CUUdDQ3NHQVFVRkJ3TUNCZ2dyQmdFRkJRY0RBVEFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFDRXhNYXVQUGczbjk1S1RsMEMrYVNoNzF5dHoxUHFXaExDdGJLU2s4T2EzRzhlMWU2OAoxVitNWVZDYjUxa25IT0YzRE03WkJBWllUYXB0MXRXelN6VnBCN2dkd2dNdU1Od3l0aDg2aXJ2eXFpbm96Q1FnCitBdHBsbDVkV0RsRVQzNFhpT05tTVVlYnVCd0JxZnNUeFArUE9IN2c0aEZpV1pTZEhRWTM4aFZFMWF5Wk92cjkKb2F1WXhpSVdEd3pzTkhCRXZaeTltVllNYnQ1dy9BYjZUWllQdTRiSVRIRE5mdlN6Y2h4aGVYNmt4VjRQQ3ByOQpSYmp0V1VDa0hqaWdJQ3AxclZhTDZGS0ZRUUxtUjFhZkxzRmxXZ1NTTzdsM1RHQUNINWZPSjR2cnBnOUtBTUk4ClNCUyt5TWsrRHRsSHljRG1QT3FBb2t3clBpUmpNVXF4N2NLeAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==`
	// token := `ZXlKaGJHY2lPaUpTVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUprWldaaGRXeDBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5elpXTnlaWFF1Ym1GdFpTSTZJbVJsWm1GMWJIUXRkRzlyWlc0dGFEazJOamNpTENKcmRXSmxjbTVsZEdWekxtbHZMM05sY25acFkyVmhZMk52ZFc1MEwzTmxjblpwWTJVdFlXTmpiM1Z1ZEM1dVlXMWxJam9pWkdWbVlYVnNkQ0lzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG5WcFpDSTZJbU5rTnpVMk5HRmhMV1kxTTJJdE1URmxOeTA1WkRoa0xUQTRNREF5TnpaaE1URTNaQ0lzSW5OMVlpSTZJbk41YzNSbGJUcHpaWEoyYVdObFlXTmpiM1Z1ZERwa1pXWmhkV3gwT21SbFptRjFiSFFpZlEuQ05OcDRRaTVQcFFLQVJuRkw1MGh5M3RWdm5vVVM2YWdwOVBwUGxIYXBmeThURGY4cEhBZm1IazIzTFRGUjlJU3k2TWVEY01VNmtJVFlNOVF3akFOYm9TNVNXekw5ZTRxNjM4NHFtdUN6Z3JVbnFfbEY2OWNyNUF5ejYtMHU0RnQ0aV8yWVVFU2xxb3ltT3o0NGg4cENtcXZ1MWRMQkdFazRMa2VhTzNucUVYWG00dHRoOE1zLTlpY294elk0WDdXZmxtMjFYV1JZbDh0anhzb2JLejVSbUdmeUNIaXNqcVF5WHpnQXRNdTRSbGVXWGR1U2tIRmJIRTZ3LU1aUGRqbWE4RmcyREFoTlJJbFc5U3BVckstT2dUZ1R6VDZCaUZSTmctU3pnU0x4U3pLWmRlX0FaRk41dl9nTS1QMHpQWFl3N0c5Q3d0X3hXemJlZ3FvbVRCYXVn`

	certif := `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNE1EUXlNakV6TWpVeU1sb1hEVEk0TURReE9URXpNalV5TWxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS2NiCjA4eXdDZHlBQ2pQSnNHK1ZVYTFRODJWWTNjTStkc09TZW1BOEx3VHQ4ckl0bVllRm1ISjRDQmpjaElkbFQ4UGoKNTJzUmhZS1E1T2JyWDVqSEVwRXR5WmtWZ1ZmM1ROMTJHUmxmUFVSNldZRzlFNzBId0VGK0o5TG9oMTV1Ny9IVQpPWk9XOGYwbnl6a2xIZjZoWU9aZ0xndWhlZ2prL2htbzdldVJRbzUwS0VBNS9LRG8yUVp4VU4rWTQwOURpZ0ZiCjNhb3o0N0VyaUhFVFRYZnRHQkQ2MURSb0I5WmdxdVBzVWY3SG9EUzJiWnVORTVrczJaUVMySVFoWmZFaElxK1UKZ0cxZGNpUE5NQzJGalNjVlJnd0VzczEyVVJXNVRPUU12b3JlR285SmFDbWNJNUUxaVRRR3plNjlpODNaUzg2Vgo3ZmhEczhCZWxvVlNyWmwrK0lNQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFFSXR0cGFLNUxCTHVCbDNIREhFWndsbnJkNncKVnBqRllWWmdZcEVnSDJxNDRhSWVITlFDR2dVWmZkK3Ewdmp6b3pxNlExLzAzTUtuOHlSZXpsQSszL09HSDJTcgpKVWpwNlp6bi9lYmovb1NSLzMwUmtBQ1ZDc0tCdHduWDBUUndGajRUTGhrVmQ1aXplVUc4Z1NVQ0w5c1Fjc1B4CjR1aGVmTjlnTnhNcEQrYWFwVVRsVUhwZ1RmR0VqR1VaQXlSQUJ1WFR4aEFCbjlRUmR2dUNPd1FOek42VmFpdlcKb2F0MVVvdFFBV3A0NWgraEJVanhXQ1ArNGlFNlM1cXFoSlZYZ05KZ0hvMHVkeU95dGNOWVphM2xiTGlPK2JYdApCQ2ZYeTVoMjgrd3dCcDRkTk5qd3hkWVBBMFJQNmpIN2htc0lBZzl5enF4bkNWTjBibW52TUp3WEJuRT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`

	// tokenka := `ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklpSjkuZXlKcGMzTWlPaUpyZFdKbGNtNWxkR1Z6TDNObGNuWnBZMlZoWTJOdmRXNTBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5dVlXMWxjM0JoWTJVaU9pSmtaV1poZFd4MElpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WldOeVpYUXVibUZ0WlNJNkltUmxabUYxYkhRdGRHOXJaVzR0WW14emR6VWlMQ0pyZFdKbGNtNWxkR1Z6TG1sdkwzTmxjblpwWTJWaFkyTnZkVzUwTDNObGNuWnBZMlV0WVdOamIzVnVkQzV1WVcxbElqb2laR1ZtWVhWc2RDSXNJbXQxWW1WeWJtVjBaWE11YVc4dmMyVnlkbWxqWldGalkyOTFiblF2YzJWeWRtbGpaUzFoWTJOdmRXNTBMblZwWkNJNkltTTNaVGd5WW1FekxUUTJNekF0TVRGbE9DMWlOalV3TFRBd01HTXlPVFpoTTJJMFppSXNJbk4xWWlJNkluTjVjM1JsYlRwelpYSjJhV05sWVdOamIzVnVkRHBrWldaaGRXeDBPbVJsWm1GMWJIUWlmUS5rNWJ5bkxJRUxITkRMQXdEUGE3ekpWYzJSRDR4eUVuTEViakNHU1R2MHZJRFhQTFZXVHlGZ0tTTjNMMENlN090RFpFTmg4dEtpWGVHUGNmRTdyWk1ydmVRVEh6dWVpcDloc0wtekItNXhzU2JKMzdfR0xERnRlQUtwNm9QckdmeWlqMENGVW9OWEd3ZV85azE3SlQ1dVBOWGJSZjJpUlFOdUJGNldRcWpkN1J6Tl90Q3JELUx4c1RBVEk4OHVRUzdqei02cnlaSEtGVjB1aEtsT3FQMVVsd1ZKTk5Kc2VhNGhTQTJ5WGpZbTgzU044UTZONDZ5T3lJeTh6VExXMlVOd3BrYWlHUmJDMTAzYkc5aXUwSE05RjhUZ240bTlQTHlVRXV5ajAwclJZS09ER3V4UlJGTVB0VEUtQTU2TGNXS2FIME93OUZhUkc3XzRKNXhrUWYxYkE=`

	// clientkey := `LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBdnQ0dy9KczdXUmlwSTRub0hpOU1wUlBlSVpFdEZTTW5TQ1o3MWRwK2hRZ0lnaW5XCjlHWUtmNStGNGkrYlZqUC9MYmcwYkpjQVhkcHdmb09sd3hVVDZ4ZDVFTHFPMHlqK2VNZmRLT20zTHJNdnVzSUQKeDltSjV2Z0JEUDZ0SVRrZjljRHJPd3FJRWpoekhRcXRqTE0yQk1HenJZNk9Ud091alpBdm9HUXZFV0pMNC9CZwp5YlhjdUtnaWluNzZGVGtleUhRZU9wcmNkVHJhSzZ3QUNaZms4amhuOS9hQWdIZTJMSE1LdTZMRWtxWDA1T2ErCkZUOFBQQ2c0QUFmcnNnNUY3L2ZmRUFlZFFSeWFpRTNqWVVKZElaenJGM2FzZFVVSXoxK2U2bk51Y0hSWFVuczEKNUJBRUxtYkJvcmlpSFdZUVJMSnZGNVlMTlVUWnpPUmZLTEwzM3dJREFRQUJBb0lCQUdLVE1HdVY1RjNNZmJZcwpDQ1JtdTBmYWJmT2FIZFYzMVZiUEFVL2VTMDk3YWFHSDZFdEsxQXM3b1JMREVZL0F4UDZnenZweU5pOUNuS3pLCll2YlEwUHV1b01rQ0FMZVB5WFVwaTlBUWZKbnkweWk2QU9mYk12eUZnMDFwenBLRkJUdVFDaXp3OEh1d2ljc3EKODV6aUJYa0piVG1xa1ZhL2lRdjF0cE00aXBLSDh0MkgvbHJaaXF1NUNrcW9PbXdmL09FMWFjWkxhWWF4dG1xTQpTM3BvS2hjZGgxWmhYNVNuZTRxQWxnTTZRWUZ1U1VqYW9PbW51QUNMM2FHNnlVM3ZHQU9xaUhrV3FTQkNMczloCkFKZHIwYmtING82QkRqSTlhTXdRcHdGSmw5YU9WOUZkdkxIZGVOR2llUTcrbzFjckJFYzRmTHlwSndLSUE5dHYKVFUzMklYRUNnWUVBK0NCdUh1WW4vdzUwb2d5aldYWkl5ZlBGWStUT1R5YWR1SG5uNUpRQ2V0Unc4b1p0Yk90bwpHUXNGRVdHeFcrOUJTdFlKSWFlMDEvdittQWhCVk1YamlYL2NHOXVwZjNkRWdvaHV4b21GRGo3dGwvOFVWMmd1CnM2Y3FOTXpnd29rc0tkYXIzMjI3RXVWUHh3NHgwc2cybGcrWndFTEg4WGJWVGozMk1FZnE4NWtDZ1lFQXhPeWoKMFRlbkpZUCtLaFRma3k5OUFQNzlnRlZTTlA3WEhxYmtXWnhFRFA0NFFiQUJFRFZ4ZzNtNWJFOUErbTBTK1VEOQpRaDVCVEptcGNDOWhDTzJ4RmdOK1FwOEtuekVzY05mVEVtYm5wYkU2L294UDMwTnBNSjg5SVBKNnZXT1YrUE03CmlFc1JVQVJTR05oMnVqdHh3S0x0VUgrenR0VGVaR0JTdjZuQzhqY0NnWUVBaHJnYzhqdm1sVzVFMTBOallZeCsKZ3VBUFdXaCt0Nnp3ejV1bzA0dWxPUW1sZFppVlN5RVplUmRwbmdGYjZkMmlwcjVGWVBlTWtnUnBQQ1NuVEI3Ugpwdk04RUFnWkpITWVTSDFKSUJURW9ISjhVQjJYN3NsTEtoSG1NWnJYb2VnV2lYVGNCc2l1WE5rU2tySmJUT1dWCjlhM3N2ZDNFYjQ4a3k0R0s3TFh2bEdrQ2dZRUFrUzZaMC95QTJYTEhwc1MrMUdlMWRFK0tHOXhMZ0ZERnpvNWkKV2dLUVZUZnp4OUgzNXJoUUdRdGIvaE1zSjdUVXdUajl2b3BKd0N5bHM5VHFhRWU5UUNxUklwTFlwT2IvQ2E3RQpxWk4rZ3pUMzlvVUJ1ZXVjR01HOXNwV3lrZ0JpcUNqRElrZWQydTFralhiQmlhbWJ3dGNidVRaOUMzVkRCS1BUClBnVHRlZDhDZ1lFQXVmSW03WHJkaFFOTEE5elhoNlBIUTc5Z1NuOHBiZXY3OVM1MTJFNS9VOWY5TEllcjJlWHcKNWp3L3dCMWFialBSVHBYNmkrYWtrUGlNbnUwS3JDZ0dJdGw0R3JMMzBvNTh2M0lUV0ZKL1h6WU94MTJPbXRrbgpzYTMzWU5kMFZnSkt0cEI4UG8vVjVKOThDN2xGQVVISGtPSEJkMHJVaFpHNk5kbVVmWW12YkhBPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`

	// apitoken := `eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImVsYXN0aWNzaGlmdC10b2tlbi02c2JzbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJlbGFzdGljc2hpZnQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJmNTg4OWMzNy01MWU4LTExZTgtOWZjYS0wMDBjMjk2YTNiNGYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDplbGFzdGljc2hpZnQifQ.OBkPonfFwK6oov0eLPyxpIWacLzC2DtYONfiGXeTetKvS-aVhZucYPQNwzVUki4TxWes3ndlipt5-OmKYTS8e79klUJHnq6YLC19K8gmnbwMIsM6OfjvUnpRSYXu03ib_8pyDfGXKs8Ntd4C9hYC22vpSihGka5KNFmh9l6m-dpuU0mYDGwljFscu2P09EX2g3NgnBzpLsVeoHbcA7mziDjDLYnArmcqf8JXdJp3uhvINo9CsAZcdIop6snfhEWeGJYeIZdp-KxJaKVi6NVH3RbJ5vwuazVH7xFUsiou_9KEbtscF9utW_ZL1ue3SmPIAm9NXDGxny64NPabcN73_g`
	apitoken := `eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InNoaWZ0bWstdG9rZW4tOTl2OG0iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoic2hpZnRtayIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImUwMmI1NGVlLTU5OWUtMTFlOC05OGMzLTA4MDAyNzZhMTE3ZCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OnNoaWZ0bWsifQ.BLo0u8zDA_xGNwY1lY8AYNwAtpNtqIyu1YRFukw0Q89ahF-szKi3BzzbS0oEA0H9RHeUEU-LWNhp4ZzBdxFxyBRkJwzJyNxLzg4EbWj8N1lNMf3F78FYmNEAQPPtJ_a5PgLLOmEtPYcmw0XuS_TIqfV2cJHv0tXO-3HEykH-caqjrXaPKSCsujhdetOxSglnlY97olAM0nvnRWPIQipOML3W1-Ij76eSUbFWbCdMpwcchDWRkOmICLgPe1Hsbbqyteh454Jy5YShvqteCcx7zOluLhPk1yDO0iFRWwEXigdbUwtEfZd-YY7yh6ywMzYqsiOWOcFy7oncEkm7wsX1ug`

	// ctx := "kubernetes-admin@kubernetes"
	// ctx := "elasticshift"
	// namespace := "elasticshift"

	ctx := "shiftmk"
	namespace := ctx
	config := clientcmdapi.NewConfig()
	config.Clusters[ctx] = &clientcmdapi.Cluster{
		Server: host,
		CertificateAuthorityData: []byte(certif),
	}

	config.AuthInfos[ctx] = &clientcmdapi.AuthInfo{
		// ClientCertificateData: []byte(certif),
		// ClientKeyData:         []byte(clientkey),
		// User:  "minikube",
		Token: apitoken,
		// Username: "kubernetes-admin",
	}
	config.Contexts[ctx] = &clientcmdapi.Context{
		Cluster:   ctx,
		AuthInfo:  ctx,
		Namespace: namespace,
	}
	config.CurrentContext = ctx

	clientBuilder := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{
		ClusterInfo: clientcmdapi.Cluster{
			InsecureSkipTLSVerify: true,
		},
	})

	//kubeconfig := filepath.Join(
	//	os.Getenv("HOME"), ".kube", "config",
	//)
	//clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// fmt.Printf("\nClient config %#v", clientConfig)

	// fmt.Println("Token = ", clientConfig.BearerToken)

	clientConfig, err := clientBuilder.ClientConfig()
	if err != nil {
		panic(fmt.Errorf("Failed to connect to kubernetes : %v", err))
	}

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(namespace)

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	/*fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	fmt.Println("Fetching logs...")

	req := clientset.CoreV1().RESTClient().Get().
		Namespace("shiftmk").
		Resource("pods").
		Name("demo-deployment-857588f755-z8lcv").
		// SubResource("log")
		// Name("demo-deployment").
		// Param("selector", "app=demo").
		// Resource("pods").
		SubResource("log").
		Param("follow", strconv.FormatBool(true))
		// Param("container", "web").
		//Param("previous", strconv.FormatBool(false))

	readCloser, err := req.Stream()
	if err != nil {
		panic(err)
	}

	fmt.Println("Streaming logs...")
	io.Copy(os.Stdout, readCloser)
}
