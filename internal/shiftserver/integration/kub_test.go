package integration

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCreateContainer(t *testing.T) {

	// host := "https://192.168.99.100:8443"
	// host := "https://10.10.7.22:6443"
	// cacert := `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQVRBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwdGFXNXAKYTNWaVpVTkJNQjRYRFRFNE1ERXdPVEV5TlRFeE5Wb1hEVEk0TURFd056RXlOVEV4TlZvd0ZURVRNQkVHQTFVRQpBeE1LYldsdWFXdDFZbVZEUVRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTTYwCjd2SjNISTh3bG9XSWF2UkJneDRzWXROZGNLeFpiRlBWYk8rajRzemVVRlN1TllCUitXdmlHMzgyelQvKzVBNDgKT28wUHhib3g3djdRQXBua1RLTGNrT21xR05DN3RlOHNRVWtFemRSWGM3Y0x0Nmh4dHBrb2NraEtwK0ltM1BJQwpvQ2hzQmMwTWxWSmpiR25IRHJXZjFZa2JxaHQ1c3lBRTBEbVVGbXJ6MWdvdEgwN0d1Rkw2b0oyWjJCeXdmaERzCjh1cFB3djR6VzJLbXV4MjRwNUNnOEhVeTV4N2NmeFhQN2tsbkY5bVFRSDVqcXp6eDRBTnBWWURJTzA0QXZJK3cKQ1JXdGtlT3hpSmxXRkRFcjF0NURtcmp2YUlYcWdGN0ROYWlhMWtpcm94czFKUlFGM0xTK21va09iNmt2N1ZNVgpMa1dwZ0VMMG9DRnZybFMyT2k4Q0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUIwR0ExVWRKUVFXCk1CUUdDQ3NHQVFVRkJ3TUNCZ2dyQmdFRkJRY0RBVEFQQmdOVkhSTUJBZjhFQlRBREFRSC9NQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFDRXhNYXVQUGczbjk1S1RsMEMrYVNoNzF5dHoxUHFXaExDdGJLU2s4T2EzRzhlMWU2OAoxVitNWVZDYjUxa25IT0YzRE03WkJBWllUYXB0MXRXelN6VnBCN2dkd2dNdU1Od3l0aDg2aXJ2eXFpbm96Q1FnCitBdHBsbDVkV0RsRVQzNFhpT05tTVVlYnVCd0JxZnNUeFArUE9IN2c0aEZpV1pTZEhRWTM4aFZFMWF5Wk92cjkKb2F1WXhpSVdEd3pzTkhCRXZaeTltVllNYnQ1dy9BYjZUWllQdTRiSVRIRE5mdlN6Y2h4aGVYNmt4VjRQQ3ByOQpSYmp0V1VDa0hqaWdJQ3AxclZhTDZGS0ZRUUxtUjFhZkxzRmxXZ1NTTzdsM1RHQUNINWZPSjR2cnBnOUtBTUk4ClNCUyt5TWsrRHRsSHljRG1QT3FBb2t3clBpUmpNVXF4N2NLeAotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==`
	// token := `ZXlKaGJHY2lPaUpTVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUprWldaaGRXeDBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5elpXTnlaWFF1Ym1GdFpTSTZJbVJsWm1GMWJIUXRkRzlyWlc0dGFEazJOamNpTENKcmRXSmxjbTVsZEdWekxtbHZMM05sY25acFkyVmhZMk52ZFc1MEwzTmxjblpwWTJVdFlXTmpiM1Z1ZEM1dVlXMWxJam9pWkdWbVlYVnNkQ0lzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG5WcFpDSTZJbU5rTnpVMk5HRmhMV1kxTTJJdE1URmxOeTA1WkRoa0xUQTRNREF5TnpaaE1URTNaQ0lzSW5OMVlpSTZJbk41YzNSbGJUcHpaWEoyYVdObFlXTmpiM1Z1ZERwa1pXWmhkV3gwT21SbFptRjFiSFFpZlEuQ05OcDRRaTVQcFFLQVJuRkw1MGh5M3RWdm5vVVM2YWdwOVBwUGxIYXBmeThURGY4cEhBZm1IazIzTFRGUjlJU3k2TWVEY01VNmtJVFlNOVF3akFOYm9TNVNXekw5ZTRxNjM4NHFtdUN6Z3JVbnFfbEY2OWNyNUF5ejYtMHU0RnQ0aV8yWVVFU2xxb3ltT3o0NGg4cENtcXZ1MWRMQkdFazRMa2VhTzNucUVYWG00dHRoOE1zLTlpY294elk0WDdXZmxtMjFYV1JZbDh0anhzb2JLejVSbUdmeUNIaXNqcVF5WHpnQXRNdTRSbGVXWGR1U2tIRmJIRTZ3LU1aUGRqbWE4RmcyREFoTlJJbFc5U3BVckstT2dUZ1R6VDZCaUZSTmctU3pnU0x4U3pLWmRlX0FaRk41dl9nTS1QMHpQWFl3N0c5Q3d0X3hXemJlZ3FvbVRCYXVn`

	// certif := `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNE1EUXlNakV6TWpVeU1sb1hEVEk0TURReE9URXpNalV5TWxvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBS2NiCjA4eXdDZHlBQ2pQSnNHK1ZVYTFRODJWWTNjTStkc09TZW1BOEx3VHQ4ckl0bVllRm1ISjRDQmpjaElkbFQ4UGoKNTJzUmhZS1E1T2JyWDVqSEVwRXR5WmtWZ1ZmM1ROMTJHUmxmUFVSNldZRzlFNzBId0VGK0o5TG9oMTV1Ny9IVQpPWk9XOGYwbnl6a2xIZjZoWU9aZ0xndWhlZ2prL2htbzdldVJRbzUwS0VBNS9LRG8yUVp4VU4rWTQwOURpZ0ZiCjNhb3o0N0VyaUhFVFRYZnRHQkQ2MURSb0I5WmdxdVBzVWY3SG9EUzJiWnVORTVrczJaUVMySVFoWmZFaElxK1UKZ0cxZGNpUE5NQzJGalNjVlJnd0VzczEyVVJXNVRPUU12b3JlR285SmFDbWNJNUUxaVRRR3plNjlpODNaUzg2Vgo3ZmhEczhCZWxvVlNyWmwrK0lNQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFFSXR0cGFLNUxCTHVCbDNIREhFWndsbnJkNncKVnBqRllWWmdZcEVnSDJxNDRhSWVITlFDR2dVWmZkK3Ewdmp6b3pxNlExLzAzTUtuOHlSZXpsQSszL09HSDJTcgpKVWpwNlp6bi9lYmovb1NSLzMwUmtBQ1ZDc0tCdHduWDBUUndGajRUTGhrVmQ1aXplVUc4Z1NVQ0w5c1Fjc1B4CjR1aGVmTjlnTnhNcEQrYWFwVVRsVUhwZ1RmR0VqR1VaQXlSQUJ1WFR4aEFCbjlRUmR2dUNPd1FOek42VmFpdlcKb2F0MVVvdFFBV3A0NWgraEJVanhXQ1ArNGlFNlM1cXFoSlZYZ05KZ0hvMHVkeU95dGNOWVphM2xiTGlPK2JYdApCQ2ZYeTVoMjgrd3dCcDRkTk5qd3hkWVBBMFJQNmpIN2htc0lBZzl5enF4bkNWTjBibW52TUp3WEJuRT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=`

	// tokenka := `ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklpSjkuZXlKcGMzTWlPaUpyZFdKbGNtNWxkR1Z6TDNObGNuWnBZMlZoWTJOdmRXNTBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5dVlXMWxjM0JoWTJVaU9pSmtaV1poZFd4MElpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WldOeVpYUXVibUZ0WlNJNkltUmxabUYxYkhRdGRHOXJaVzR0WW14emR6VWlMQ0pyZFdKbGNtNWxkR1Z6TG1sdkwzTmxjblpwWTJWaFkyTnZkVzUwTDNObGNuWnBZMlV0WVdOamIzVnVkQzV1WVcxbElqb2laR1ZtWVhWc2RDSXNJbXQxWW1WeWJtVjBaWE11YVc4dmMyVnlkbWxqWldGalkyOTFiblF2YzJWeWRtbGpaUzFoWTJOdmRXNTBMblZwWkNJNkltTTNaVGd5WW1FekxUUTJNekF0TVRGbE9DMWlOalV3TFRBd01HTXlPVFpoTTJJMFppSXNJbk4xWWlJNkluTjVjM1JsYlRwelpYSjJhV05sWVdOamIzVnVkRHBrWldaaGRXeDBPbVJsWm1GMWJIUWlmUS5rNWJ5bkxJRUxITkRMQXdEUGE3ekpWYzJSRDR4eUVuTEViakNHU1R2MHZJRFhQTFZXVHlGZ0tTTjNMMENlN090RFpFTmg4dEtpWGVHUGNmRTdyWk1ydmVRVEh6dWVpcDloc0wtekItNXhzU2JKMzdfR0xERnRlQUtwNm9QckdmeWlqMENGVW9OWEd3ZV85azE3SlQ1dVBOWGJSZjJpUlFOdUJGNldRcWpkN1J6Tl90Q3JELUx4c1RBVEk4OHVRUzdqei02cnlaSEtGVjB1aEtsT3FQMVVsd1ZKTk5Kc2VhNGhTQTJ5WGpZbTgzU044UTZONDZ5T3lJeTh6VExXMlVOd3BrYWlHUmJDMTAzYkc5aXUwSE05RjhUZ240bTlQTHlVRXV5ajAwclJZS09ER3V4UlJGTVB0VEUtQTU2TGNXS2FIME93OUZhUkc3XzRKNXhrUWYxYkE=`

	// clientkey := `LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcFFJQkFBS0NBUUVBdnQ0dy9KczdXUmlwSTRub0hpOU1wUlBlSVpFdEZTTW5TQ1o3MWRwK2hRZ0lnaW5XCjlHWUtmNStGNGkrYlZqUC9MYmcwYkpjQVhkcHdmb09sd3hVVDZ4ZDVFTHFPMHlqK2VNZmRLT20zTHJNdnVzSUQKeDltSjV2Z0JEUDZ0SVRrZjljRHJPd3FJRWpoekhRcXRqTE0yQk1HenJZNk9Ud091alpBdm9HUXZFV0pMNC9CZwp5YlhjdUtnaWluNzZGVGtleUhRZU9wcmNkVHJhSzZ3QUNaZms4amhuOS9hQWdIZTJMSE1LdTZMRWtxWDA1T2ErCkZUOFBQQ2c0QUFmcnNnNUY3L2ZmRUFlZFFSeWFpRTNqWVVKZElaenJGM2FzZFVVSXoxK2U2bk51Y0hSWFVuczEKNUJBRUxtYkJvcmlpSFdZUVJMSnZGNVlMTlVUWnpPUmZLTEwzM3dJREFRQUJBb0lCQUdLVE1HdVY1RjNNZmJZcwpDQ1JtdTBmYWJmT2FIZFYzMVZiUEFVL2VTMDk3YWFHSDZFdEsxQXM3b1JMREVZL0F4UDZnenZweU5pOUNuS3pLCll2YlEwUHV1b01rQ0FMZVB5WFVwaTlBUWZKbnkweWk2QU9mYk12eUZnMDFwenBLRkJUdVFDaXp3OEh1d2ljc3EKODV6aUJYa0piVG1xa1ZhL2lRdjF0cE00aXBLSDh0MkgvbHJaaXF1NUNrcW9PbXdmL09FMWFjWkxhWWF4dG1xTQpTM3BvS2hjZGgxWmhYNVNuZTRxQWxnTTZRWUZ1U1VqYW9PbW51QUNMM2FHNnlVM3ZHQU9xaUhrV3FTQkNMczloCkFKZHIwYmtING82QkRqSTlhTXdRcHdGSmw5YU9WOUZkdkxIZGVOR2llUTcrbzFjckJFYzRmTHlwSndLSUE5dHYKVFUzMklYRUNnWUVBK0NCdUh1WW4vdzUwb2d5aldYWkl5ZlBGWStUT1R5YWR1SG5uNUpRQ2V0Unc4b1p0Yk90bwpHUXNGRVdHeFcrOUJTdFlKSWFlMDEvdittQWhCVk1YamlYL2NHOXVwZjNkRWdvaHV4b21GRGo3dGwvOFVWMmd1CnM2Y3FOTXpnd29rc0tkYXIzMjI3RXVWUHh3NHgwc2cybGcrWndFTEg4WGJWVGozMk1FZnE4NWtDZ1lFQXhPeWoKMFRlbkpZUCtLaFRma3k5OUFQNzlnRlZTTlA3WEhxYmtXWnhFRFA0NFFiQUJFRFZ4ZzNtNWJFOUErbTBTK1VEOQpRaDVCVEptcGNDOWhDTzJ4RmdOK1FwOEtuekVzY05mVEVtYm5wYkU2L294UDMwTnBNSjg5SVBKNnZXT1YrUE03CmlFc1JVQVJTR05oMnVqdHh3S0x0VUgrenR0VGVaR0JTdjZuQzhqY0NnWUVBaHJnYzhqdm1sVzVFMTBOallZeCsKZ3VBUFdXaCt0Nnp3ejV1bzA0dWxPUW1sZFppVlN5RVplUmRwbmdGYjZkMmlwcjVGWVBlTWtnUnBQQ1NuVEI3Ugpwdk04RUFnWkpITWVTSDFKSUJURW9ISjhVQjJYN3NsTEtoSG1NWnJYb2VnV2lYVGNCc2l1WE5rU2tySmJUT1dWCjlhM3N2ZDNFYjQ4a3k0R0s3TFh2bEdrQ2dZRUFrUzZaMC95QTJYTEhwc1MrMUdlMWRFK0tHOXhMZ0ZERnpvNWkKV2dLUVZUZnp4OUgzNXJoUUdRdGIvaE1zSjdUVXdUajl2b3BKd0N5bHM5VHFhRWU5UUNxUklwTFlwT2IvQ2E3RQpxWk4rZ3pUMzlvVUJ1ZXVjR01HOXNwV3lrZ0JpcUNqRElrZWQydTFralhiQmlhbWJ3dGNidVRaOUMzVkRCS1BUClBnVHRlZDhDZ1lFQXVmSW03WHJkaFFOTEE5elhoNlBIUTc5Z1NuOHBiZXY3OVM1MTJFNS9VOWY5TEllcjJlWHcKNWp3L3dCMWFialBSVHBYNmkrYWtrUGlNbnUwS3JDZ0dJdGw0R3JMMzBvNTh2M0lUV0ZKL1h6WU94MTJPbXRrbgpzYTMzWU5kMFZnSkt0cEI4UG8vVjVKOThDN2xGQVVISGtPSEJkMHJVaFpHNk5kbVVmWW12YkhBPQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=`

	// apitoken := `eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImVsYXN0aWNzaGlmdC10b2tlbi02c2JzbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJlbGFzdGljc2hpZnQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJmNTg4OWMzNy01MWU4LTExZTgtOWZjYS0wMDBjMjk2YTNiNGYiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDplbGFzdGljc2hpZnQifQ.OBkPonfFwK6oov0eLPyxpIWacLzC2DtYONfiGXeTetKvS-aVhZucYPQNwzVUki4TxWes3ndlipt5-OmKYTS8e79klUJHnq6YLC19K8gmnbwMIsM6OfjvUnpRSYXu03ib_8pyDfGXKs8Ntd4C9hYC22vpSihGka5KNFmh9l6m-dpuU0mYDGwljFscu2P09EX2g3NgnBzpLsVeoHbcA7mziDjDLYnArmcqf8JXdJp3uhvINo9CsAZcdIop6snfhEWeGJYeIZdp-KxJaKVi6NVH3RbJ5vwuazVH7xFUsiou_9KEbtscF9utW_ZL1ue3SmPIAm9NXDGxny64NPabcN73_g`
	// apitoken := `eyJhbGciOiJSUzI1NiIsImtpZCI6IiJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InNoaWZ0bWstdG9rZW4tOTl2OG0iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoic2hpZnRtayIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6ImUwMmI1NGVlLTU5OWUtMTFlOC05OGMzLTA4MDAyNzZhMTE3ZCIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OnNoaWZ0bWsifQ.BLo0u8zDA_xGNwY1lY8AYNwAtpNtqIyu1YRFukw0Q89ahF-szKi3BzzbS0oEA0H9RHeUEU-LWNhp4ZzBdxFxyBRkJwzJyNxLzg4EbWj8N1lNMf3F78FYmNEAQPPtJ_a5PgLLOmEtPYcmw0XuS_TIqfV2cJHv0tXO-3HEykH-caqjrXaPKSCsujhdetOxSglnlY97olAM0nvnRWPIQipOML3W1-Ij76eSUbFWbCdMpwcchDWRkOmICLgPe1Hsbbqyteh454Jy5YShvqteCcx7zOluLhPk1yDO0iFRWwEXigdbUwtEfZd-YY7yh6ywMzYqsiOWOcFy7oncEkm7wsX1ug`

	//encoded_kubefile := "ewogICAgImtpbmQiOiAiQ29uZmlnIiwKICAgICJhcGlWZXJzaW9uIjogInYxIiwKICAgICJwcmVmZXJlbmNlcyI6IHt9LAogICAgImNsdXN0ZXJzIjogWwogICAgICAgIHsKICAgICAgICAgICAgIm5hbWUiOiAia3ViZWNsdXN0ZXIiLAogICAgICAgICAgICAiY2x1c3RlciI6IHsKICAgICAgICAgICAgICAgICJzZXJ2ZXIiOiAiaHR0cHM6Ly8xMzAuMTk4LjY2LjM0OjI2ODU1IiwKICAgICAgICAgICAgICAgICJjZXJ0aWZpY2F0ZS1hdXRob3JpdHktZGF0YSI6ICJMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VaU1ZFTkRRWGt5WjBGM1NVSkJaMGxLUVV0T09VODJNblJuUzJKeFRVRXdSME5UY1VkVFNXSXpSRkZGUWtOM1ZVRk5SR3Q0VG5wQk1VSm5UbFlLUWtGTlRVeHFhM3BaVkZsNlRWZFpOVTlVWnpGT2VsSm9XbXBaTlZreVdUUk9ha3BvV21wWk1FNTZTbXBPYWtrMFRGZDBNVmx0Vm5saWJWWXdXbGhOZEFwWk1rVjNTR2hqVGsxVVozZE9ha0V4VFdwTk1FMTZUVEpYYUdOT1RrUlZlRTFFU1hoTmFrMHdUWHBOTWxkcVFUVk5WR04zVGxGWlJGWlJVVVJFUXpRMUNrMHlSVEpOZWtadFQxUnJORTVVWXpCWlYxa3lUMWRPYlU5RVdYbFpWMWt5VGtSamVWbDZXWGxQUXpGeVpGZEtiR050Tld4a1IxWjZURmRPYUUxSlNVTUtTV3BCVGtKbmEzRm9hMmxIT1hjd1FrRlJSVVpCUVU5RFFXYzRRVTFKU1VORFowdERRV2RGUVRVMVJXVTRRM2xFY0RSdFVVMUJOek5DY1ZGUU5UQm9Vd280Wmpoc1NuZFdVamhVVERsak9GTnFVRlJrZWpsTVFYSkJaSE5uZEROUFVHWTFUVEJGTm1aeWRtVjRlRVZFZHl0eWVuQndjM05QSzJKTllsRjROalJ6Q2tsc1dtTmFOV3d3YldjNFUwdFhTVW8zVVZCNlpXeENUa3hCVW1ZM1FUSmFaMmhPYm5SVVVrdHFkVTR5UkdSMVFUZFhkbFpGTlZReE9EQkJXa0YyVURNS1JrZzVSMU4zTUdOTlNHbHhLM0lyWldKYVZqWkdaV1V4YkVOTVdrcFdNVlZEYW5aUFZtaEZNMHRtTkVWcVlXeGxVVEY0VmpGVVUwWjBlWE5UVVhwVk5RcFNNMnhFTWpObGJUZHRhbXQwUjFKU1IxRXZhMGhGTDA5VlRVVjVUalJPVFZaMGVYRnlOMkU1TVZobVRtSmxTa05SVG01c2VIVnNObUphZGtvdmFVTnVDalU1UXpSMWJqZHdTVnBQYlVNMFJVbDJOSFVyTkdGU2FqZHdVSFIwTWt0VU1HVlJlVnBSUkdnMEwxWTRWakJsSzJsUGFESlpOVUYyYVZSaU5saHdaakFLWkRZMFZuUXpTemQzWlhsSE5XRkdaRFpxZVRCWmNuTk1kVTAwTlROUlFXcHFUVUZUYm5ObVVGRnJkWEJNT1doVk5FNUljWGhDTDFKRWQzTnJXVWRPU0FwUE1VSjJjVmd5VGxoMmNrWlNURTVJY0hGclYyUjZSbEpYUzBSd1lYRnJWR0prYW1WTVJWUnZOVU56VDJoUU5qZ3JkMmhST1RSaGNYZHROMmN5VGtkYUNrWjNWbGM1YWpCd1NHeFhkbEl5TkdabGVtdEVMMFZPS3poMmNFcHlhVlF4TldvMGNsbDZlbHB3ZUZKT1UyUnBSbHByUkhGWlYwTkhhbkpKU0hKVGMxb0tVSEJGTDJoM09FTXJXUzl4ZFRkSFFrYzVVU3M1Y1RGelZGZDJkazFTWVRGeFdqTlNOV05XYjJGMVVVNVZRV2QxYW05VlltNTZTV1ZXY0RaV1ZrSkxid3BDTkZoUkwwVmhUVkZrU2xwcmRYRnpPWEpSZGpWTU1GSlNRemt3TlVoVVNIZ3pTMDkyVlhCQ1ZHaFhUM1U1ZUVWQmVUZFBhVmhtUVdWMGRrUXlORFJSQ2xSMmEzcDJWR2ROTWtZdk5EZEZXaTlNY0d0RFFYZEZRVUZoVGxGTlJUUjNTRkZaUkZaU01FOUNRbGxGUms4eGNFRnRjeXQ0VWtwS1VFOWhVbUY1YVRFS1JuVnFNRlJJZG5oTlFqaEhRVEZWWkVsM1VWbE5RbUZCUms4eGNFRnRjeXQ0VWtwS1VFOWhVbUY1YVRGR2RXb3dWRWgyZUUxQmQwZEJNVlZrUlhkUlJncE5RVTFDUVdZNGQwUlJXVXBMYjFwSmFIWmpUa0ZSUlV4Q1VVRkVaMmRKUWtGSWVYRkJRVUpVWjNKNGEzZGpNMGwxZUV4ckt6RlRTVGxxTnpORlYzSnJDbk5IYlhBMWVGZG9Ta2h2ZUZSU2JGcHVTVTlEYVZad2Iza3lUVnBFZHpkRlJXSlliakZIUzFvMFdGbHhTMGd5YzI5QlEzcElka1p1VG1scWVISnRWbGtLUjBJNGVHNWFXR1JHVVVNeWR6QTViSFJpYTFGb2NtNVhWRUpOTDJwTVMyMDBlbTAyWXpkdmVUbEtNekZpY1VwdFVucGxiSE5sZDNGU2FVZHhZblZJU0FwNmVGRnZaM2t3V21aUk5WZHNabVk1TDNadFJXUkdkSFp4WWl0SlJFTXdTVGRUV0ZkV1pFaG1XVlpDYUdFdlNuTnlTbnBhWTFGblJIWXZhMEZwVUVWVENtYzFXak42SzNsUWFIUXJTRXR0UTBkS1VFTjRZMWRTY3psS1prazFUVVkzYVhoc1FuaGFRbHBvY2sxcFMwOU9TVVowUjFORWEzcHlWSEV6Y2tGUE1td0tUbFExY2pBNFkxUTBiV2hKUkhSV01GWXZUVVpPWkc5Q09HMXdNRXR6TkZJeE5uSkpTbWh0YlZOR1pqZEhUVFV6VEhFMlFtazFjVTlNVEZOcGRGa3JSZ3BTU0dsNWVESjFPWGx0V0VWeWJuWnpWRFpQUVV0bWVsQXlablUyUkZKbEsxRlRkQ3RUY0d4NVVIZFFWamhwVGxsSU5tNXBlRGswWVZoa1YxaENPVm93Q21vd01tZHFVVkpGTURSbmMycEhWazFsWjNNNFJXdEVhelZEZEUxdFVUazRORWRIUTNReFIwbEVZblZTT0hCek1IVmFiMjlpYzJ4UVdrUkVMM2RRYVhNS1ZYUlVVR0YzTmtoUVkzRTNWVWR5U0V3MFVUQkxTRmRSWmxSdWNUVk1UazR3V0hoNGRGVkpUMDlsU20xRmMweEdhMGxyUlhkbVYxQndTak16Wlhoa1RBb3hNMlpTVDBJeU9GQnJVRWRMVTNSMFVqazBNbXBWWnpGdmFIQnpOSEJySzNSd1NuUlpVV0Y1ZHpJeVpFcGxRMUZsTWpFek5UZHpOWGRzVEVSeVV6Tk1Da2cxYmtveVVubzBhRWxHWm5ScFRuVnlTREIzYzBwWk1UbDJVMVZLY2xOWlYzSjJWR1IyWlZWMmJuVlBPVmh3VFhOcFVGcHBSak5SU0RsdFUwWnJSVVFLVEZaRmFXSldXbXh0VkZCTkNpMHRMUzB0UlU1RUlFTkZVbFJKUmtsRFFWUkZMUzB0TFMwSyIKICAgICAgICAgICAgfQogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAibmFtZSI6ICJsYWJjbHVzdGVyIiwKICAgICAgICAgICAgImNsdXN0ZXIiOiB7CiAgICAgICAgICAgICAgICAic2VydmVyIjogImh0dHBzOi8vMTAuMTAuNy4xNjE6NjQ0MyIsCiAgICAgICAgICAgICAgICAiY2VydGlmaWNhdGUtYXV0aG9yaXR5LWRhdGEiOiAiTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVTjVSRU5EUVdKRFowRjNTVUpCWjBsQ1FVUkJUa0puYTNGb2EybEhPWGN3UWtGUmMwWkJSRUZXVFZKTmQwVlJXVVJXVVZGRVJYZHdjbVJYU213S1kyMDFiR1JIVm5wTlFqUllSRlJGTkUxRVZYcE5SRUUwVGtSTk1VNXNiMWhFVkVrMFRVUlZlVTU2UVRST1JFMHhUbXh2ZDBaVVJWUk5Ra1ZIUVRGVlJRcEJlRTFMWVROV2FWcFlTblZhV0ZKc1kzcERRMEZUU1hkRVVWbEtTMjlhU1doMlkwNUJVVVZDUWxGQlJHZG5SVkJCUkVORFFWRnZRMmRuUlVKQlNsbHNDbXhUWm5KV2JrTkRaQzlZTVRWU1kxUlVObGRSTTNkc2JtZDRkVGxPT1RaS1kzaFZhMHRSVUhCbFlXYzBaV3RPTXk5Qk1EWjBhM0V2V0cwck5tbDZWRFFLWmpocGJub3dUVTlIZFZObE5reFdTV2haY0VRdlVEQTNVazVLY0dOMlRUTkthRGRuT1VGTVJtTjBjR2hGZWt4VE1HNHZORTkyTDJFMFdtVjFWbWgxVVFwVGFYSXlNV3RoYUdOVlYwOW1TamR2VFU1VWJHbGhWM2wzVXpNdlNXNXVOR2h0WVZOVmVVcEtjMEZFUlZjeVN6SnBlbTFrY1RSdk0yNVFhWGwzVVNzMENtVklSek5FY2pKWlJrSk5XVlkxTmsxck0yNVRWR2hTWWxBMlFYRlBVRXBtTkZkUVRYRjVSRk5uVEVoWFdFaGhkV0puUW5kdFRtSlFaamhHY21KUGEyUUtSVzE0U2tWalJpODJTbFphWmtoMFJ6VTVaVVJLTm1OemRGcGFWWGxEWlV4R1FUSmxSVUpVTjJNeGNVSmFVWFpWZGtjNVMxRldjbEIzVm1wM1l6Uldad3BsVFhsa05tUlNWR3g0U2toa05FSjFiMnRyUTBGM1JVRkJZVTFxVFVORmQwUm5XVVJXVWpCUVFWRklMMEpCVVVSQlowdHJUVUU0UjBFeFZXUkZkMFZDQ2k5M1VVWk5RVTFDUVdZNGQwUlJXVXBMYjFwSmFIWmpUa0ZSUlV4Q1VVRkVaMmRGUWtGSlMwaEJNWGxoUVVwSmIzSnpObVZyVWxsMFVsWkdlaTlMUmtnS1ZVdGlNa3RtWjFCQ09XVjBWRUZhZWxJdmFISXhhVlpDYURZdlZ5OXVRblJVTmtnMmFVWkpWSFJuVWtkbVduWklZalZCVlVkdGVETnNUbGxvZFVWR1J3cFNSVW8wY0hsNGIzVXdaamhwTW1Wbk0zZzBOa3hqVUM5bVNuQTVNVzB4TVVkVllrUXpTMlpZZEdOT1ZYaFZLMUV4T0VWU2RWaFBWVVJvYTJKeFJFaFRDbWhaWWk5VFYxcFdhMWtyTUZCQ1RFdFZkRXBUUjA5SGNXNVBlamxSZVdvMFRtbFlSRVJ5YTBSRVFXRnJWVXhNY25GS1RXZzNMMjFsVVhST01YVlVXamNLYzJOSlpIUjNVMVZUZDJabmFtOUpiM2d3WjJ4UlltaFBURk41WlZoM1QyVnBURGd2TkZsUlptb3llRVZsV0hwWlFtNTBLM2w0TjAwM2FVTXdRbE5sVlFwbldsTlNPWGNyUmtkUFoySktlWE5ZVTNKWUwwc3lMM012U1ZFM2RsQnhiRXRLUnk5blNWcEdRV2xqVTJobk5GRmhaVlpJT1VsUFpsUm5hejBLTFMwdExTMUZUa1FnUTBWU1ZFbEdTVU5CVkVVdExTMHRMUW89IgogICAgICAgICAgICB9CiAgICAgICAgfSwKICAgICAgICB7CiAgICAgICAgICAgICJuYW1lIjogIm1pbmlrdWJlIiwKICAgICAgICAgICAgImNsdXN0ZXIiOiB7CiAgICAgICAgICAgICAgICAic2VydmVyIjogImh0dHBzOi8vMTkyLjE2OC45OS4xMDA6ODQ0MyIsCiAgICAgICAgICAgICAgICAiY2VydGlmaWNhdGUtYXV0aG9yaXR5LWRhdGEiOiAiTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVTTFla05EUVdNclowRjNTVUpCWjBsQ1FWUkJUa0puYTNGb2EybEhPWGN3UWtGUmMwWkJSRUZXVFZKTmQwVlJXVVJXVVZGRVJYZHdkR0ZYTlhBS1lUTldhVnBWVGtKTlFqUllSRlJGTkUxRVkzZE9ha1V6VFhwamQwNHhiMWhFVkVrMFRVUmpkMDVFUlROTmVtTjNUakZ2ZDBaVVJWUk5Ra1ZIUVRGVlJRcEJlRTFMWWxkc2RXRlhkREZaYlZaRVVWUkRRMEZUU1hkRVVWbEtTMjlhU1doMlkwNUJVVVZDUWxGQlJHZG5SVkJCUkVORFFWRnZRMmRuUlVKQlRVZHZDbEF5YlVvMWNIVkVjMVZFY0hZMlVUTlBlVWxUYjJFek5qbHpPR1o0Vm14eVRWUkViemxaZVV4TmFIcExSRzk1Tkc1WE1reHJhRzFEU1RGT0szcFNNVWtLYVhoQ1QzVkthbTloVVhFcmJrWnVkMU5SYkZkQ1UxQkpjamhOUmtac1ZFdE9TM3BhVjFCd1NVRlhZMUZ3YUhKV2RIZzNRMVpsUTI1dlIzVnljMU5hTXdwU1RrMDVaM0JRUmtSSlJVcEpOV1JuYld0M2FsQlRTV2RJUTFZeVZrSk5lREpJZUZwamRWRjBSRVpyZFRSS1JYZFFTbHBKUW5OcE1EVTRXV1ZUWnpkbkNsbGhNRFpUV0VjM1MzUTFRVmhOU3k4NGMzSjROMFpMTm1sTlMyMW1UMmQ1WW5SRVdXcEVPV28yY0hGTFkzUmFXRzlZWkVkQlprcEtSbWh5U2tFMVRsY0tjMGhsYWsxUWFWZDNhMnRSSzFoTlFrRmFRWFZhV0hOUU5FdDNXa0pzWTAxQ0wySmtlRWxrWlV0clRYcERjbmRCTnpOSmVFZFVNRUphV1N0R1RGbFdRd3BYVTJ0eU9XNW5RM2xUT0Zjd1VTdEVSa0Z6UTBGM1JVRkJZVTVEVFVWQmQwUm5XVVJXVWpCUVFWRklMMEpCVVVSQlowdHJUVUl3UjBFeFZXUktVVkZYQ2sxQ1VVZERRM05IUVZGVlJrSjNUVU5DWjJkeVFtZEZSa0pSWTBSQlZFRlFRbWRPVmtoU1RVSkJaamhGUWxSQlJFRlJTQzlOUVRCSFExTnhSMU5KWWpNS1JGRkZRa04zVlVGQk5FbENRVkZCWlRVeVRVODRURzltUW5FM1pUaDVOVlIzUmpaQmVDOHpUa3BzTlc5Tk9EbFpZWGRrYURsellUTXJOU3RZZFRGSlZ3bzJNRlUwTlhWVlUzTkhTVEpMUlhabllTdHBXalIwTUZkemQzbFRTM2gxVUM5dVNGaE1Za1pWVGpoNlJVOTFlRVpQYjJOU1NVSm9VbWg2WlVFeloxTlpDbHBTVmt4RmRFVmFRV2hrZWprd1NFcDBMM0JLUW01MGFETlZOV2RSU2taaWQxRm9VRGs0U2tKMlRUQmhSazFhTjNKVmJ6aG1iVU5LTnpkNlNFTXdTVUlLUlU5RVJFcERlRmhuTkU5RE1WUmxiamxvWTFjclFrbE5NbUZPZWpCMmVrSk5kVXBrY1ZsbVIyeEpaVkZuTmxJeGMyeFdaRGc1YlZJeE0wWmFjVGRJYndwTlRXNWpOV3MwWTNWSVFtMWpSMUJ2YTNZeU5sZDZlSHBoVkhSNk5FWlBUVmhNVEVvdk1tb3dUU3MwVkd4UlJsQjZSVWM1YnpkVWNFOHdXVWxYYlhBd0NuUnFUVmx0WW1GT09ITnphMEppZEVGSmFTdDNiRGgyYkhOcmR6Z3ZWR3RTV2tzNVN3b3RMUzB0TFVWT1JDQkRSVkpVU1VaSlEwRlVSUzB0TFMwdENnPT0iCiAgICAgICAgICAgIH0KICAgICAgICB9CiAgICBdLAogICAgInVzZXJzIjogWwogICAgICAgIHsKICAgICAgICAgICAgIm5hbWUiOiAiZ2hhejR1QGdtYWlsLmNvbSIsCiAgICAgICAgICAgICJ1c2VyIjogewogICAgICAgICAgICAgICAgImF1dGgtcHJvdmlkZXIiOiB7CiAgICAgICAgICAgICAgICAgICAgIm5hbWUiOiAib2lkYyIsCiAgICAgICAgICAgICAgICAgICAgImNvbmZpZyI6IHsKICAgICAgICAgICAgICAgICAgICAgICAgImNsaWVudC1pZCI6ICJieCIsCiAgICAgICAgICAgICAgICAgICAgICAgICJjbGllbnQtc2VjcmV0IjogImJ4IiwKICAgICAgICAgICAgICAgICAgICAgICAgImlkLXRva2VuIjogImV5SnJhV1FpT2lJeU1ERTNNVEF6TUMwd01Eb3dNRG93TUNJc0ltRnNaeUk2SWxKVE1qVTJJbjAuZXlKcFlXMWZhV1FpT2lKSlFrMXBaQzAxTlRBd01EQXdNVVJFSWl3aWFYTnpJam9pYUhSMGNITTZMeTlwWVcwdWJtY3VZbXgxWlcxcGVDNXVaWFF2YTNWaVpYSnVaWFJsY3lJc0luTjFZaUk2SW1kb1lYbzBkVUJuYldGcGJDNWpiMjBpTENKaGRXUWlPaUppZUNJc0ltVjRjQ0k2TVRVek1EYzRNRFExT1N3aWMyTnZjR1VpT2lKcFltMGdiM0JsYm1sa0lpd2lhV0YwSWpveE5UTXdOemMyT0RVNWZRLmY1d0QtbEx0S2dWWUs3MkVwWEtRbUZQZUpEZnQ0amN6d0dRbG1YeFpLUXNLazREdW9SVTg3V09kVzZoaTFOQ2U5bjFtbHByQmVrY2otM0NtbnQ0aXhZWFBYMWFWX2VFU05NOGs1LUNRZ0p5UFBqSWtLQjY1QUk3N1U3YnV1UHM3T1JhUlV4T2ZIVzNoMk14NV9TaE9CRkM5QlNqM1ZYY0RMcW53aXhfempBOWEwYUFZaDdHM2lZN0hJaEJzVWxEb0RocWNDV3hZc3VsQkl4MHhoY1NOcVdjS2tSbFAwOFJUTE5KOGtWYlBRcGtldnFJbHlWVUp5RWxBV05lcTRtMXkzNDZOYjZoUG01aF8zSTJHbVZGX0Y5WUdQM1RHS21oZE1IS3hPSVM0MEM4TnFqLXNzYmpNWl9TeS1mTGd0akM1el9lS1NvR002S25fZDYyczhzT0g2ZyIsCiAgICAgICAgICAgICAgICAgICAgICAgICJpZHAtaXNzdWVyLXVybCI6ICJodHRwczovL2lhbS5uZy5ibHVlbWl4Lm5ldC9rdWJlcm5ldGVzIiwKICAgICAgICAgICAgICAgICAgICAgICAgInJlZnJlc2gtdG9rZW4iOiAiSjFCNHZmVWpLT2V5OWNPY1ZLS3lwOTB0RlFwZE5adXBzU05OQkZlYkZSdTRVUmlkdFBOVWRhTS1LYTdRSXlNTTNrbGYtQy01dmoxeHBLTTZCVW9zd0ZubVZEaGpvUjMwbGYtSGwwYU5yU0xQcDM5YXNHZmdMNjVsUC05RVg0WlhwbVQzM0ZuMHUtZzBWbmlEWkl3WWxzYkZOcnhORXAyWmRoZFJBbzROdHBrVS1DZUN1OGhKdlM2UHlpT2JxOUhvendrZ1VCSzJTcUliYjdCUUNHZmFKWlBfVk1wRlNKYnd5WUtsdE5UV1BHNnFnTkFITGlSRElmRXlWLTV4bGFydkRkeGppUS1vUnU5dk5xdlpSTGY4UFgwZnlYR3VrTkVSTmdHczR6elN3M3BER0h0RzV6LWZ1UWtMczBWZnY1cHppdzhBN0U3bTN1QWVUX19DZEdqX2dTRlQ3bkJncl9CQ1B1NWlKb2tmM044a2hEVFYtdGZIbm43dDBqNFM5S0dIZDIxT195ZWNmSWlBMEZVOHhDWl9LVzJ3aHVGOXFpZEpTbGk0OVlUUkEtUW9UWDJhWllfa0lxRkcyaTRKMUw4UTduRlNpRXA3aW5vMVB6dThDUkgwdHZRcHY2RnBWZjBOb1RaOGNrbGVpNm81bDctZGJ5cEtMbG5Ia1Zrd1BoV0ljLUhwVzZMUUwyQXh4NnFtdDlqcWQxWDFjVU1IRExvNDdja0VIMDdFLXRaek5ueFBGZW83WGRJemlhc2ZEUnFfOEttN0piVnJnSU9QTG1DYXQ0U0sxT2tZVWhqYnZMZjB3dXp3VWlUZ1BpelRTa3lGalNkSlJfcEF2ZTVlaktEWFVJTHFmeFB6c09OTF9MbzdOdnFOYk1pWkxrTkJkOGhXeWJ6Ujh4MlVVYk0xOUNQQ3Q5Q0hFeE0wQXVYbjlxNkw0eFc1OFJWMC1Mb0NpZ1o1RkJnWFo1V01vblRJYlYwLWIyNUJDekRLMEpSNjVsalBvWHlCa3VLZDkyQldRREExWFFKX0VGZ1FqRXhGbTltU0xoVXJfdFo1NGZoWWQ1djBGMUJYSjdWQ2FrY1JKQy1vSG1vdVJpdG5hOG1lcGhEaEoyakhuNkJ3WjFHZ2hCVF8yR1Y0Ymw3S3BuQ1RMYnM5Rkw0RVBCUlpHMS04b3ciCiAgICAgICAgICAgICAgICAgICAgfQogICAgICAgICAgICAgICAgfQogICAgICAgICAgICB9CiAgICAgICAgfSwKICAgICAgICB7CiAgICAgICAgICAgICJuYW1lIjogImt1YmVybmV0ZXMtYWRtaW4iLAogICAgICAgICAgICAidXNlciI6IHsKICAgICAgICAgICAgICAgICJjbGllbnQtY2VydGlmaWNhdGUtZGF0YSI6ICJMUzB0TFMxQ1JVZEpUaUJEUlZKVVNVWkpRMEZVUlMwdExTMHRDazFKU1VNNGFrTkRRV1J4WjBGM1NVSkJaMGxKWkdob2JVUklRMmwzV1c5M1JGRlpTa3R2V2tsb2RtTk9RVkZGVEVKUlFYZEdWRVZVVFVKRlIwRXhWVVVLUVhoTlMyRXpWbWxhV0VwMVdsaFNiR042UVdWR2R6QjRUMFJCTVUxNlFYZFBSRkY2VGxSYVlVWjNNSGhQVkVFeFRYcEJkMDlFVVRCTlJFWmhUVVJSZUFwR2VrRldRbWRPVmtKQmIxUkViazQxWXpOU2JHSlVjSFJaV0U0d1dsaEtlazFTYTNkR2QxbEVWbEZSUkVWNFFuSmtWMHBzWTIwMWJHUkhWbnBNVjBackNtSlhiSFZOU1VsQ1NXcEJUa0puYTNGb2EybEhPWGN3UWtGUlJVWkJRVTlEUVZFNFFVMUpTVUpEWjB0RFFWRkZRVzlTYVdReFlsQjNUemRSUm10V2N5c0tjVk5qTW5aeWJuWllZV0pDZVRkMkwwUk9NRXBxTW1wU1Z6SnRjR1JDTmtaU1lWVnpRVkl2ZURWaU1rSmxTV00zTkdkVWNYQlhTM1Z4VDNkd1dHVlFaUXBYWTFsa2MyMWpiVU5hWW5oNFlVVmtZMUJQVFVaNFJTODRibVZ1UTNrNWQyRjZhRGRrYlhKcVZrMVJWRXBTYUZwdVIzUXZObGRHV2twbFFVaElSM1Z0Q2tkTU9IbEthVk5FTVZNeE1WVnVRVTF6TURGQmEwUllOakphUnpJME5HdHJaSGxQUlhKek1FaDJZVWhPUWtGRWVsVjJNRmR0VmswMVNuWTRNMllyVlVvS04zZFFVbFYxYUVWT1pHbGhlRFoxUjBGT1VITm5WM0Z4VWpsWFp5czRMMU5tS3paWWN6Wm1VemxNYUhGQlRVaHlhV1E1TDJ4MmJrMTZLMDF0UkdGMVNncDFPWEY0YzBrMU1ESkJXVlpLV21GNFNXSkZiVEJXZFVZellVTTNWR0l4YWpWSVFqVlNkVXBzV1daUFpEQkVRVzlFYkhSQlRVWnVNM0pWUW5oUVJGbExDakJyZGpWemQwbEVRVkZCUW05NVkzZEtWRUZQUW1kT1ZraFJPRUpCWmpoRlFrRk5RMEpoUVhkRmQxbEVWbEl3YkVKQmQzZERaMWxKUzNkWlFrSlJWVWdLUVhkSmQwUlJXVXBMYjFwSmFIWmpUa0ZSUlV4Q1VVRkVaMmRGUWtGRVdESTNZVkoxUmxOWGVWaE1iWHB6WVVGSWJuZDNZbXBVVTNabU9XOTRkWEpZWkFwMFR6bDRNRGhpVlZONWJtZ3JTR3ROVUdJMk1YbEpOR0pOUkdSNlYycEdNMlpUVG5CTE4zaFRaa1ZqUkVJclNuaEVVWFJ5TjFGVkwydHpXVmhvYjB0Q0NrZDBkVGRNVjNnclJGSkhZVU5YTDFkRlJERlVlVTl3Y0ZWb00xazRkRk5PVGpWeU1uSmxaMmt6WjB0cU4xTk1VM1Z3Y210Q1NtTjRNRU0yUTNWd2RUVUtNelpFUlZjMGVWZEtaa1IzU0VOT1VtOXJZbloxTmtsaFpsaFNXa1E1U2tsR1dHWlNVSFo0UTJjckwwWXZXWGM1WTFCa1VVUm5WVXhEVkdJeWVrRlFWd3BZYnpGcFUxTndPR2RpTVc1VmNEUTFRME14TkhOMmJVRlpkVU5QUWxSYVpETmpjRmxWTmtsWVdpOVRPRnBwWTBZMlNWVTFTMmQ1ZUVseFNUWk1WbEpzQ2tkSVNHaGhLMGN4YUUxNk9UQTJlbWRCVGs1NE1XSmpjWFJ3Wms4eGJEUldTbkV3VUZOd1NuVTFjbUYwUkVNMU1HMUJVVDBLTFMwdExTMUZUa1FnUTBWU1ZFbEdTVU5CVkVVdExTMHRMUW89IiwKICAgICAgICAgICAgICAgICJjbGllbnQta2V5LWRhdGEiOiAiTFMwdExTMUNSVWRKVGlCU1UwRWdVRkpKVmtGVVJTQkxSVmt0TFMwdExRcE5TVWxGYjNkSlFrRkJTME5CVVVWQmIxSnBaREZpVUhkUE4xRkdhMVp6SzNGVFl6SjJjbTUyV0dGaVFuazNkaTlFVGpCS2FqSnFVbGN5YlhCa1FqWkdDbEpoVlhOQlVpOTROV0l5UW1WSll6YzBaMVJ4Y0ZkTGRYRlBkM0JZWlZCbFYyTlpaSE50WTIxRFdtSjRlR0ZGWkdOUVQwMUdlRVV2T0c1bGJrTjVPWGNLWVhwb04yUnRjbXBXVFZGVVNsSm9XbTVIZEM4MlYwWmFTbVZCU0VoSGRXMUhURGg1U21sVFJERlRNVEZWYmtGTmN6QXhRV3RFV0RZeVdrY3lORFJyYXdwa2VVOUZjbk13U0haaFNFNUNRVVI2VlhZd1YyMVdUVFZLZGpnelppdFZTamQzVUZKVmRXaEZUbVJwWVhnMmRVZEJUbEJ6WjFkeGNWSTVWMmNyT0M5VENtWXJObGh6Tm1aVE9VeG9jVUZOU0hKcFpEa3ZiSFp1VFhvclRXMUVZWFZLZFRseGVITkpOVEF5UVZsV1NscGhlRWxpUlcwd1ZuVkdNMkZETjFSaU1Xb0tOVWhDTlZKMVNteFpaazlrTUVSQmIwUnNkRUZOUm00emNsVkNlRkJFV1Vzd2EzWTFjM2RKUkVGUlFVSkJiMGxDUVVOTlpXRmlTV1ozZG5oUk9UTklSQXB6V0hCUmRWTmxSa0V3UTI5VVRrVnpWRFpMZGxGU2JFcFpXRXRNU3pCNlRVcHZUMVJUVnl0S1ZtdFBaVTR5TldWR1RtcHNPVE14VFdOUlZWVlVjMUpSQ2tjM1dEUmpXbkpzZWpZMFFtZFFNWHBpVTJGc1oxZEhOU3RtV0ZONlpFTnlZVGxsUkM5aE5VbERVSG9yY21GdVJtdERkVUZhSzJOd2RVeFRRVVJUZWs0S04zVjVkMHRoSzFaQlJWSnNNelZETUV0eFEyWTROMWs0Wm1zMFZteEljVkZaVEVkTVdXRm9XbkJ0YUM5a2FtSXpSRmg0ZEd0S1drSm5SSFozWVM4M09BcGxOVXhHT0M5VWVYTlNPVzVWZUVscFJreDFNR3R3WVU4eE5VeEVaMHBJZGtvd04wMWpLME5VTlROa2RWVklUa2xsUjJKcE0wWnBNWFl5ZW1SdlZYVkdDaTl4VEhoNlptRnRWMlJ5WlV0clJHVkliV1JNYUdwV2NrcGhhSFJtVW0xS09YSTRSM0pCVTJ0MU1XNUJjSE5vUjBGb2MwMHpSa2w1YUV4Q1MyUkxjVEVLUVhWMWNYRTFhME5uV1VWQmVDOHZjRGw2UTFaRkszSXJTM3BwV1hOclRIVkdUalIwSzFsNlV6ZEdUa2hoV1VNM2NERm5NM1VyUTIxb1MzZDFhRkJMZVFwVWJrRXlTRVZLYldoMUt6SkZiMUZRVUhOelVVWXZLMnBuYWk5MmJrcHFWVElyTTJjMU5FOWxLeXN6V1dwcVZWSmFNek5qV0VFeFJqaExOamxGZGtsS0NtRk9RbWd4TUVOaVlrWjJSVWQyS3pSUVlUaE1lakZ6VjNSYVoyeG5SbGhyYmtzeFNuTmxhSEo0TUVOV1ZIQjZVRmN5WnpGM1REaERaMWxGUVhwcVVWUUtPR3RYVW5wcVlrNW5VVWxCVERkR2VISXhkWFozV1hJemMzY3JOMjF6VEVGRVVVZGFZV3hvYzNkRlZrMVBSMDV1WWpOc1pUZHlNbk5pVG04cmVVVlFRd3BFTVU5TFltNTFjVE5oYjBVNVRqWldVRFpTYTNab2MzRllURzFhVUZKSFl5dDBkbEJxVWprMFZVUTVkV3hOT0dNeGRtNHhabkY1UldVNWJtSlFTekpVQ2pWdFNYQTNVelZhVjFoM1ZGQkNWbkpHTkd0dVJXSkhRMHRWWkVwNFNHNHJVMEY1U3pCQk1FTm5XVVZCYkZoWVlVVnlRbkoxUVdKV1EzUmxkM2hSWW1ZS2FGSnFSVzl6YkhGc2RFWnNaa3RMU0cxbVpVUk1lRkZxV0U5TFVHOVRjamwwVlVOSFRETm5kbmh2WXpsUFdHdFlTelJZTURoUWNFbzFkbmxCZUVSWmVRcFlSMnBaZFdseVUyeE9TMkZGZVhZeFdtSklXa00zUjJadmN6TjVNMnRYYVZsdlNGaHVSbFJpTUZrMlEySjFja3hrVFRGeWRHWnZUV0p2YWxVdllVMTVDbmhUVGpZdk1HdENkMnhYZUZCd05WSkdOalEwTDFkRlEyZFpRVTlwTmk5dmFFTXpVRE5TTDB3dlQyTnJaR2xFZEROTmNtVklWRFY0TTBOa1FYVkpTWGdLY1VzME5IbFViMHR4YUdwQ1ZsbHVSRmhGTVhSREsxbDNNalJvV1U5cFpHcFRNa3BXWm5SWVJIZDVNazFXU0Uwd2VGSlVlVTFDY0dnMWFHNUJURUpFYkFwNFZVZDZNRGhrTlRaaEsxVlJZbEJ2UjJ0UVQwSnlPV3gxVXpFd056QXlhV280VVdwWVoxTm9aMnhaWkVFcmJrWkNSVVJGVWpoQ1JuaDVhVllyYjBWTkNtUm5ORWwwVVV0Q1owSkllbnBFVVV0YVlVTkxWV281UVRGbVpsbGlURlY2YkdsU1FteE5hSEZvZFhCU2FHVnpaWFl5Ynk5MlpXcEdjRVZNWkVSa2NWWUtkMHByZVN0aFJtOWhMME16VDBWVVdYWnZRWGxPZDNONVRYVm5lbVpYVFdKSlRHaGxhV0pOZVdSRWJITmxla1JDT0dwQ1Izb3pabTR4VW1OclNrZEVOQXBoUTA1UWVuSnFXbG8xZVZwTU0zSnBOM1JsVWs5b1IzRnJabTlyVEU1M01taFBUemMzUlRrM1NEUldWVVZHWVVSd2RIZFpDaTB0TFMwdFJVNUVJRkpUUVNCUVVrbFdRVlJGSUV0RldTMHRMUzB0Q2c9PSIKICAgICAgICAgICAgfQogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAibmFtZSI6ICJtaW5pa3ViZSIsCiAgICAgICAgICAgICJ1c2VyIjogewogICAgICAgICAgICAgICAgImNsaWVudC1jZXJ0aWZpY2F0ZS1kYXRhIjogIkxTMHRMUzFDUlVkSlRpQkRSVkpVU1VaSlEwRlVSUzB0TFMwdENrMUpTVVJCUkVORFFXVnBaMEYzU1VKQlowbENRV3BCVGtKbmEzRm9hMmxIT1hjd1FrRlJjMFpCUkVGV1RWSk5kMFZSV1VSV1VWRkVSWGR3ZEdGWE5YQUtZVE5XYVZwVlRrSk5RalJZUkZSRk5FMUVZM2RPZWtFelRsUmplVTFXYjFoRVZFVTFUVVJqZDA5RVFUTk9WR041VFZadmQwMVVSVmhOUWxWSFFURlZSUXBEYUUxUFl6TnNlbVJIVm5SUGJURm9Zek5TYkdOdVRYaEdha0ZWUW1kT1ZrSkJUVlJFVnpGd1ltMXNjbVJYU214TVdGWjZXbGhKZDJkblJXbE5RVEJIQ2tOVGNVZFRTV0l6UkZGRlFrRlJWVUZCTkVsQ1JIZEJkMmRuUlV0QmIwbENRVkZFUzBOWWJ6UXlNMUEzYlZWbmVFSk1RMVo0ZEVwaWNuRnJVWGN2TnpJS1ZqazJhQzlKZEZwaGNuaFlXbXQwVUZkQmFraHdMMjQzWmtoc2FqbDFabll6VEcxWWFFSnZNVkpNT1ZWVmRIRkJkemRuT0hOWWREZGxValJKVEdaWVVRcHRhRTh3ZUdVMWNVWm1UR1JLYWtKTFdFRnVjelUwVnpWTVFuSkZaamRZV0VvMFZuQlpNVzAwU0RCaGEyOTVWbGRXT1hWd1JYUm1jbU1yU21aRE9GcFhDbVpzWW14TE9XNW5VRk5XVDNwNVQxbEtaVXczUkN0MWJXUkZVMGxSY3pKeVZsUXJWbk5EYnprNWEwMDJMMmhEYlVsUE4xWTNkR3hZZHpWc1QxSXpOVFFLVDBoSE1UZDZOMm96VjFjNGMzVkRRVTFWU1d4a2RYRkJhRWMyTkUxbVpXTTNkbTFFWlZod2JWQkxla1ozYkd0RGMwcFdjelptZW1aWU9XSk5jbWR0YXdwMmJITkxhMll4VXpCMlIzSnNabVZuWm5oblJHa3JUMk5hTUVkbUsyRTBUbWhZVGpsNWEwdzRTVlUyUldkYVJuVk1lWE5HUmt4dFZrRm5UVUpCUVVkcUNsQjZRVGxOUVRSSFFURlZaRVIzUlVJdmQxRkZRWGRKUm05RVFXUkNaMDVXU0ZOVlJVWnFRVlZDWjJkeVFtZEZSa0pSWTBSQlVWbEpTM2RaUWtKUlZVZ0tRWGRKZDBSQldVUldVakJVUVZGSUwwSkJTWGRCUkVGT1FtZHJjV2hyYVVjNWR6QkNRVkZ6UmtGQlQwTkJVVVZCY2tvemVFaDRlbEEzTTA5b01VOXBNUXB3YjAxdGJWSlNWR1JzYzNoU1ltdFpiVXBuUjBsQk1HbFFXbGhUYUVJcmRYUjROemRrVjNOR1kzaDRZV2swVFhvMU4waDFlalU1U1VaT2JtcDBTMkp5Q21ReGFGTnhWMmM0ZUVoaWFVdHFWVEJtVWtNMVZXOXNWWEJ6VUdNeU5VNXlaeXR2Yld4NlRtSXZTMjlwVkhNelZEZ3ZNRTVtVVhaNWN5dGxNelpuWmpnS2NrTk9Mell2YjA5M2NVMXhWell4TTFSUlNWaEhka2xZWjJRNUt6QlJXVTF5Wms5WFJuWTBNV05CVW1WQ1oxTnZXbVJoYVZveFRXbFRibUV5SzFkNWJRcDJkMjA1YUc5c1VqUlRNblJVWmpCYVZFMVhhRFptUmpkWGRVRkhWR0k1UVZoTFZFMHdXV2xLWVZrMlkxbzBWMGxFY1dSM1FqY3JUVFJEWWpsc2NrOTVDazFLWldKbE5Xa3dVV3BDWjBwRU0yaHRabGcxUmtsTFZYaG5Zak1yT0dwdVZrcHdURmsyUTNsMGJEVldiM2hvUXpFMFpXZ3lhbWxuZFdwc1ZVMDNkWFlLVldoNWVsRkJQVDBLTFMwdExTMUZUa1FnUTBWU1ZFbEdTVU5CVkVVdExTMHRMUW89IiwKICAgICAgICAgICAgICAgICJjbGllbnQta2V5LWRhdGEiOiAiTFMwdExTMUNSVWRKVGlCU1UwRWdVRkpKVmtGVVJTQkxSVmt0TFMwdExRcE5TVWxGY0ZGSlFrRkJTME5CVVVWQmVXZHNOazlPZEhvck5XeEpUVkZUZDJ4allsTlhOalp3UlUxUUt6bHNabVZ2Wm5sTVYxZHhPRll5V2t4VU1XZEpDbmcyWmpVck0zZzFXUzlpYmpjNWVUVnNORkZoVGxWVEwxWkdUR0ZuVFU4MFVFeEdOMlV6YTJWRFF6TXhNRXB2VkhSTldIVmhhRmg1TTFOWmQxTnNkMG9LTjA5bFJuVlRkMkY0U0NzeE1YbGxSbUZYVGxwMVFqbEhjRXROYkZac1ptSnhVa3hZTmpOUWFWaDNka2RXYmpWWE5WTjJXalJFTUd4VWN6aHFiVU5ZYVFvcmR5OXljRzVTUldsRlRFNXhNVlV2YkdKQmNWQm1Xa1JQZGpSUmNHbEVkVEZsTjFwV09FOWFWR3RrSzJWRWFIaDBaVGdyTkRreGJIWk1UR2RuUkVaRENrcFlZbkZuU1ZKMWRVUklNMjVQTnpWbk0ydzJXbXA1YzNoalNscEJja05XWWs5dU9ETXhMMWQ2U3pSS2NFdzFZa053U0RsVmRFeDRjVFZZTTI5SU9Ga0tRVFIyYW01SFpFSnVMMjExUkZsV2VtWmpjRU12UTBaUGFFbEhVbUpwT0hKQ1VsTTFiRkZKUkVGUlFVSkJiMGxDUVVOUFVGaFlja3RpV0ZOWlQzaFBaQXB1Y2sxUFFuQktjV2RKYkZvMVprNVlOMlpEVnpoaVVIZFVSM2RhY21sSGRDOW5WMlZPUTNGdFkzWnpPRkozUWpVeWFGTnZRa3hFYUdjNFpVVjFZall5Q2s1RFVrTnVUbWhVY21kNVJXZHRTMUEyS3pSRFpFZzRjR3M1UTBjeWVrMU1WbWhtYW1aRE4yMWhNazkzYTNOSWJYazVNa3g0WjNaYVVGZGlWemRGYUhnS1QwWkdaVFZ5UVdwbGJrWXJUV2NyYUU1M2QzQkdjbFZaVEdkMlluRjZibGcwTm5kTVVrdG9SRGhtU0VSbVVXNVdkRE5wUW5relRWZE5lV2RYUVVzNVpnbzNPSGRGZFZOd09XMVdXR05GZFhWcFVVNDRTVVl6TkZoUWRYWTBXVEppWm1WMWQyNXJibUZHWTJaa1drMDVaMnhwZVhrMlpVUnhUVkk0Ym1kSlpDdExDaXRqTDBoNlRtZEthVFV5Y1hKbU5XSnRZazU0VTJwdFQwRm5OMWhGZFhCM2IwbFJibFZyVURST2VqWndNbWxhU25CbGNXMW1UVGg0TkRnd01XSnhPVkFLVUdwb1dIQnRSVU5uV1VWQk9XcFJXSEF3TDFOTlUwVnJRVmhxVm5obWNGaDNXbnBuV0ZNclJ6SmxaSEZSY2toRFJGTndkVGRRU1hwSmN6aExUbmxzYmdwT1NuQnhUV3RFY25SYWFERmFVelJrV25aNlltSkRNVEJyVVZCbU5FVmxSMnhyUVhVeVJtb3lVVGRpUmxRNFdDdFZWbTAxZGtwUU1IRk9XbWRWUjNBNUNsTk1aMjFsVG1SVFZYZDFaWFpJZFd4clJqTkVaWEUyTkhaVFRtaG9LMmxYSzJkdlYwTmtkbUZYU3pSM2VsZE9NR1l3VDJndlFqQkRaMWxGUVRCb1Rqa0tkeXR3VEV4NlQxQTJSMHREUlhRMlowSXpiMUJVUVhsYU1YWlRSa1Z5TjFwU2FEUjFObVZuUm0xdFZrcFRhalYyVkRNdllsSXJOMW96T1ZOS1prWkNRd3B3VDJSMk1ERkxhMk5KVTA5clJUQTNXbEV6ZVZkV1RIUjBLMmgzZGtWamFHSkxXbGQyYVVwc2EyOVFkbFZvWjIwM1ZFVkxiVEI2TVU1aVN6bENabkpEQ2xaeU1FWXljR2xuVkhneFRGSnZWblYyTURsWWNGRlFVRlZRUWtkNEwyNVNlVTUzT1VOa2EwTm5XVVZCTVV4RFlqTnZXWFZhZUV4cFRVRk9ZbUZtVTJrS05ucG1ReXMzTW1KSFdrODFWVWhUVURGNEwxcFNRV3BZT1VSSk15dFFPV3ByVTFnd2NVRlZPVzB4ZEdoNk1EZFFVamRGV2xCYVRHNUxTQzlaV0RoQ2NRcHVSa3BIU21SQ01ISnVjWFZoWjFBMVpDOHllbGg0TTFOT2FXczRaWHAzUkdveFlrZHdXVFJsUzBjdldtRmhja2hzVEdaeWNqUm9VbEE1VUZrMVltZG5DamxXUkZKVFUwMU5kRm9yUm5WWGF6VnpkR1ZUYTBkRlEyZFpSVUZyWWpkelJtbDRTMnhUYWxWM2Iyb3hZbFp1S3psTVJDdFJLMnBMU1hCUWNXSjVUSGdLV1haQk1UVkNiMUpMVW1vd1ZFRnZZemsyTjJjeE0zTnBkRmRQYWt4NFoxQTFTa05qZVV0UmFHRmtOMDVsWkdaaVUxUlVkakp2UkhGYVFteHFabGg0Y2dvclFrVXZVbEpLWjBoa1NtcERlbGQxWWtGclIweDBRM0J1ZFU5MWRVdE9UV2xTYTBSSmNrZGFTbk53VVVjeFJXWnJTMHhtVkRaRWFFTkhjSEU0VUhGc0NsUlNSbVpXTmtWRFoxbEZRWFJ6V0VrMWNrMTJZV2RtY2tSeFRUaHJPV3RMZURSblJGbG9RaTl6ZWtwTk1VMUhNVzlDTjFkcEwwRXJRM1pPZFhkb1lsWUtVRmxPWjBZemIyb3piRmxyTkhWVk4wRlJkRmxqYURoRFIzVnpaMnBuYkRCc1JYQnNWMHBqT1dOSk0yeExhbXBHWVc0NVRXbFdVRWRLYTA4MlpVRjVWUW94ZUVkTk1YTTRTMmR3WWxvMmF6aEpNMDE2U1dsU05FcDJkWGhYZEUxUGJYVTJjMjR6Y3poNWVGTnVWV3BMYmtoaWVqaGxkRUZOUFFvdExTMHRMVVZPUkNCU1UwRWdVRkpKVmtGVVJTQkxSVmt0TFMwdExRbz0iCiAgICAgICAgICAgIH0KICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICAgIm5hbWUiOiAic2hpZnRtayIsCiAgICAgICAgICAgICJ1c2VyIjogewogICAgICAgICAgICAgICAgInRva2VuIjogImV5SmhiR2NpT2lKU1V6STFOaUlzSW10cFpDSTZJaUo5LmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUprWldaaGRXeDBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5elpXTnlaWFF1Ym1GdFpTSTZJbk5vYVdaMGJXc3RkRzlyWlc0dGJIUnVZM2dpTENKcmRXSmxjbTVsZEdWekxtbHZMM05sY25acFkyVmhZMk52ZFc1MEwzTmxjblpwWTJVdFlXTmpiM1Z1ZEM1dVlXMWxJam9pYzJocFpuUnRheUlzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG5WcFpDSTZJakF3TlRRek1XTmhMVGd6T1dRdE1URmxPQzA1WlRVekxUQTRNREF5TnpKa1pUUTJZU0lzSW5OMVlpSTZJbk41YzNSbGJUcHpaWEoyYVdObFlXTmpiM1Z1ZERwa1pXWmhkV3gwT25Ob2FXWjBiV3NpZlEuTjY3QndhNkxiN1RkT1BDMVZObzJvRTVjN2JobG5WWE1sSXprVUlaTTh2MlhSeC1jU1o1UG5vSFE4aFBsRGhqX2tZVTN3M0xSd2IyYnY1cDk1S1dtdmFXaXN0V1gwQVdzMHh2SmpEVGpPQTRTOWYyMUlaWjRvekMzemxBLS1OSTNyT0c3aXRwSU1LRHkyN3hQbWFVeU0wZzg2UEN2WUQ5VmFXNlRBNmVCZXBBN3BtcnhHeVEzejl2YTdUdmlFZDlXM2F2TTlWenBYQTJMa0xiTFQ3bDk3QjRGWlJGY2hCMU9BbTJyX2haZmFQQ0ZmUkN3RkhDbEpYM2JKdFdEM3RYSVI5YVpjYU9WV0ZJUGd4dFlNalFNaHJXUFk1UHpYM3hWUGxUeDFHdU5hcEt3VFp3dFJJUm5hLW5rcVN2XzNXYXZEeC1jaUgyUDg2cl9YTERONzk3WG1BIgogICAgICAgICAgICB9CiAgICAgICAgfQogICAgXSwKICAgICJjb250ZXh0cyI6IFsKICAgICAgICB7CiAgICAgICAgICAgICJuYW1lIjogImt1YmVjbHVzdGVyIiwKICAgICAgICAgICAgImNvbnRleHQiOiB7CiAgICAgICAgICAgICAgICAiY2x1c3RlciI6ICJrdWJlY2x1c3RlciIsCiAgICAgICAgICAgICAgICAidXNlciI6ICJnaGF6NHVAZ21haWwuY29tIiwKICAgICAgICAgICAgICAgICJuYW1lc3BhY2UiOiAiZGVmYXVsdCIKICAgICAgICAgICAgfQogICAgICAgIH0sCiAgICAgICAgewogICAgICAgICAgICAibmFtZSI6ICJsYWJjbHVzdGVyIiwKICAgICAgICAgICAgImNvbnRleHQiOiB7CiAgICAgICAgICAgICAgICAiY2x1c3RlciI6ICJsYWJjbHVzdGVyIiwKICAgICAgICAgICAgICAgICJ1c2VyIjogImt1YmVybmV0ZXMtYWRtaW4iCiAgICAgICAgICAgIH0KICAgICAgICB9LAogICAgICAgIHsKICAgICAgICAgICAgIm5hbWUiOiAibWluaWt1YmUiLAogICAgICAgICAgICAiY29udGV4dCI6IHsKICAgICAgICAgICAgICAgICJjbHVzdGVyIjogIm1pbmlrdWJlIiwKICAgICAgICAgICAgICAgICJ1c2VyIjogIm1pbmlrdWJlIgogICAgICAgICAgICB9CiAgICAgICAgfSwKICAgICAgICB7CiAgICAgICAgICAgICJuYW1lIjogInNoaWZ0bWsiLAogICAgICAgICAgICAiY29udGV4dCI6IHsKICAgICAgICAgICAgICAgICJjbHVzdGVyIjogIm1pbmlrdWJlIiwKICAgICAgICAgICAgICAgICJ1c2VyIjogInNoaWZ0bWsiLAogICAgICAgICAgICAgICAgIm5hbWVzcGFjZSI6ICJzaGlmdG1rIgogICAgICAgICAgICB9CiAgICAgICAgfQogICAgXSwKICAgICJjdXJyZW50LWNvbnRleHQiOiAibGFiY2x1c3RlciIKfQo="

	// kubefile, _ := base64.StdEncoding.DecodeString(encoded_kubefile)

	// config, err := clientcmd.Load([]byte(kubefile))
	// if err != nil {
	// 	panic(err)
	// }

	// ctx := "kubernetes-admin@kubernetes"
	// ctx := "elasticshift"
	// namespace := "elasticshift"

	// ctx := "shiftmk"
	// namespace := ctx
	// config := clientcmdapi.NewConfig()
	// config.Clusters[ctx] = &clientcmdapi.Cluster{
	// 	Server: host,
	// 	CertificateAuthorityData: []byte(certif),
	// }

	// config.AuthInfos[ctx] = &clientcmdapi.AuthInfo{
	// 	// ClientCertificateData: []byte(certif),
	// 	// ClientKeyData:         []byte(clientkey),
	// 	// User:  "minikube",
	// 	Token: apitoken,
	// 	// Username: "kubernetes-admin",
	// }
	// config.Contexts[ctx] = &clientcmdapi.Context{
	// 	Cluster:   ctx,
	// 	AuthInfo:  ctx,
	// 	Namespace: namespace,
	// }
	// config.CurrentContext = ctx

	// clientBuilder := clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{
	// 	ClusterInfo: clientcmdapi.Cluster{
	// 		InsecureSkipTLSVerify: true,
	// 	},
	// })

	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Printf("\nClient config %#v", clientConfig)

	// fmt.Println("Token = ", clientConfig.BearerToken)

	// clientConfig, err := clientBuilder.ClientConfig()
	// if err != nil {
	// 	panic(fmt.Errorf("Failed to connect to kubernetes : %v", err))
	// }

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	// deploymentsClient := clientset.AppsV1().Deployments(DefaultNamespace)

	// deployment := &appsv1.Deployment{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name: "demo-deployment",
	// 	},
	// 	Spec: appsv1.DeploymentSpec{
	// 		Replicas: int32Ptr(1),
	// 		Selector: &metav1.LabelSelector{
	// 			MatchLabels: map[string]string{
	// 				"app": "demo",
	// 			},
	// 		},
	// 		Template: apiv1.PodTemplateSpec{
	// 			ObjectMeta: metav1.ObjectMeta{
	// 				Labels: map[string]string{
	// 					"app": "demo",
	// 				},
	// 			},
	// 			Spec: apiv1.PodSpec{
	// 				Containers: []apiv1.Container{
	// 					{
	// 						Name:  "web",
	// 						Image: "nginx:1.12",
	// 						Ports: []apiv1.ContainerPort{
	// 							{
	// 								Name:          "http",
	// 								Protocol:      apiv1.ProtocolTCP,
	// 								ContainerPort: 80,
	// 							},
	// 						},
	// 						VolumeMounts: []apiv1.VolumeMount{{Name: "vol", MountPath: "/opt/elasticshift"}},
	// 					},
	// 				},
	// 				Volumes: []apiv1.Volume{
	// 					{
	// 						Name: "vol",
	// 						VolumeSource: apiv1.VolumeSource{NFS: &apiv1.NFSVolumeSource{
	// 							Server:   "10.10.7.151",
	// 							Path:     "/nfs/elasticshift",
	// 							ReadOnly: false,
	// 						}},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// Create Deployment
	fmt.Println("Creating deployment...")
	//result, err := deploymentsClient.Create(deployment)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	fmt.Println("Fetching logs...")

	/*req := clientset.CoreV1().RESTClient().Get().
		Namespace("default").
		Resource("deploy").
		Name("5b8d07f2dc294a7e86941251").
		// SubResource("log")
		// Name("demo-deployment").
		// Param("selector", "app=demo").
		// Resource("pods").
		SubResource("pods").
		// Param("selector", KW_BUILDID+"=5b8d07f2dc294a7e86941251").
		Param("follow", strconv.FormatBool(true))
	// Param("container", "web").
	//Param("previous", strconv.FormatBool(false))

	*/
	req := clientset.CoreV1().RESTClient().Get().
		Namespace("default").
		Resource("pods").
		Name("5b8e9d75dc294a5cf1b0e8f9-8486ff8988-4rntr").
		// SubResource("log")
		// Name("demo-deployment").
		//	Param("selector", KW_BUILDID+"=5b8d07f2dc294a7e86941251").
		// Resource("pods").
		SubResource("log").
		//Param("container", "f21953efa9a9a971a34980f53d0a1941e496e8e6a1a27d42e92b9eb466967872").
		Param("follow", strconv.FormatBool(true))
	//Param("previous", strconv.FormatBool(false))

	readCloser, err := req.Stream()
	if err != nil {
		panic(err)
	}

	fmt.Println("Streaming logs...")
	io.Copy(os.Stdout, readCloser)
}
