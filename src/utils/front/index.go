package front

import (
	"fmt"
	"github.com/oimoyu/OimoAdmin/src/utils/functions"
	"path"
)

func GenerateIndex() error {
	wwwRootDir := functions.GetWWWRootDir()
	indexPath := path.Join(wwwRootDir, "index.html")

	//// skip if exist
	//if functions.IsPathExists(indexPath) {
	//	return nil
	//}

	content := `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8" />
    <title>oimo admin</title>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1, maximum-scale=1"
    />
    <meta http-equiv="X-UA-Compatible" content="IE=Edge" />
    <link
      rel="stylesheet"
      title="default"
      href="./sdk/sdk.css"
    />
    <link rel="stylesheet" href="./sdk/helper.css" />
    <link
      rel="stylesheet"
      href="./sdk/iconfont.css"
    />
    <script src="./sdk/sdk.js"></script>

    <script src="https://unpkg.com/vue@2"></script>
    <script src="https://unpkg.com/history@4.10.1/umd/history.js"></script>
    <style>
      html,
      body,
      .app-wrapper {
        position: relative;
        width: 100%;
        height: 100%;
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <div id="root" class="app-wrapper"></div>
    <script>
      (function () {
        let amis = amisRequire('amis/embed');
		const amis_lib = amisRequire('amis');
        const match = amisRequire('path-to-regexp').match;

        // 如果想用 browserHistory 请切换下这处代码, 其他不用变
        // const history = History.createBrowserHistory();
        const history = History.createHashHistory();

        const app = {
          type: 'app',
          brandName: 'Oimo Admin',
          header: {
			  "type": "container",
			  "style": { "margin-left": "auto" },
			  "body": {
				"label": "Logout",
				"type": "button",
				"confirmText": "Are you sure to logout?",
				"level": "danger",
				"onEvent": {
				  "click": {
					"actions": [
					  {
						"actionType": "ajax",
						"args": {
						  "api": {
							url: "./logout",
							"method": "post",
							"data": {
							}
						  }
						},
					  },
	
					]
				  }
				}
	
			  }
			},
          // footer: '<div class="p-2 text-center bg-light">底部区域</div>',
          // asideBefore: '<div class="p-2 text-center">菜单前面区域</div>',
          // asideAfter: '<div class="p-2 text-center">菜单后面区域</div>',
          api: './site.json'
        };

        function normalizeLink(to, location = history.location) {
          to = to || '';

          if (to && to[0] === '#') {
            to = location.pathname + location.search + to;
          } else if (to && to[0] === '?') {
            to = location.pathname + to;
          }

          const idx = to.indexOf('?');
          const idx2 = to.indexOf('#');
          let pathname = ~idx
            ? to.substring(0, idx)
            : ~idx2
            ? to.substring(0, idx2)
            : to;
          let search = ~idx ? to.substring(idx, ~idx2 ? idx2 : undefined) : '';
          let hash = ~idx2 ? to.substring(idx2) : location.hash;

          if (!pathname) {
            pathname = location.pathname;
          } else if (pathname[0] != '/' && !/^https?\:\/\//.test(pathname)) {
            let relativeBase = location.pathname;
            const paths = relativeBase.split('/');
            paths.pop();
            let m;
            while ((m = /^\.\.?\//.exec(pathname))) {
              if (m[0] === '../') {
                paths.pop();
              }
              pathname = pathname.substring(m[0].length);
            }
            pathname = paths.concat(pathname).join('/');
          }

          return pathname + search + hash;
        }

        function isCurrentUrl(to, ctx) {
          if (!to) {
            return false;
          }
          const pathname = history.location.pathname;
          const link = normalizeLink(to, {
            ...location,
            pathname,
            hash: ''
          });

          if (!~link.indexOf('http') && ~link.indexOf(':')) {
            let strict = ctx && ctx.strict;
            return match(link, {
              decode: decodeURIComponent,
              strict: typeof strict !== 'undefined' ? strict : true
            })(pathname);
          }

          return decodeURI(pathname) === link;
        }

        let amisInstance = amis.embed(
          '#root',
          app,
          {
            location: history.location,
            data: {
            },
            context: {
            },
		  	locale: navigator.language.startsWith('zh') ? 'zh' : 'en-US'

          },
          {
            // watchRouteChange: fn => {
            //   return history.listen(fn);
            // },
            updateLocation: (location, replace) => {
              location = normalizeLink(location);
              if (location === 'goBack') {
                return history.goBack();
              } else if (
                (!/^https?\:\/\//.test(location) &&
                  location ===
                    history.location.pathname + history.location.search) ||
                location === history.location.href
              ) {
                // 目标地址和当前地址一样，不处理，免得重复刷新
                return;
              } else if (/^https?\:\/\//.test(location) || !history) {
                return (window.location.href = location);
              }

              history[replace ? 'replace' : 'push'](location);
            },
            jumpTo: (to, action) => {
              if (to === 'goBack') {
                return history.goBack();
              }

              to = normalizeLink(to);

              if (isCurrentUrl(to)) {
                return;
              }

              if (action && action.actionType === 'url') {
                action.blank === false
                  ? (window.location.href = to)
                  : window.open(to, '_blank');
                return;
              } else if (action && action.blank) {
                window.open(to, '_blank');
                return;
              }

              if (/^https?:\/\//.test(to)) {
                window.location.href = to;
              } else if (
                (!/^https?\:\/\//.test(to) &&
                  to === history.pathname + history.location.search) ||
                to === history.location.href
              ) {
                // do nothing
              } else {
                history.push(to);
              }
            },
            isCurrentUrl: isCurrentUrl,
            theme: 'cxd',
			enableAMISDebug: false,
			responseAdaptor(api, payload, query, request, response) {
				if (response && response.data && response.data.code === 401) {
				  const login_msg = response && response.data && response.data.o_msg ? response.data.o_msg : ''
				  amis_lib.toast.error('Not Login: ' + login_msg);
				  if (window.location.hash !== '#/login/') {
					if (window.app_data && window.app_data.is_login !== undefined) {
					  window.app_data.is_login = false;
					}
					history.push("/login/");
					//location.reload();
				  }
				}else{
					if (response.data && response.data.o_msg){
						if (response.data.code === 200) {
						  amis_lib.toast.success(response.data.o_msg);
						}else{
						  amis_lib.toast.error(response.data.o_msg);
						}
					}
				}



				return payload;
			}
          }
        );

        history.listen(state => {
          amisInstance.updateProps({
            location: state.location || state
          });
        });
      })();
    </script>
  </body>
</html>
`

	err := functions.CreateFileWithContent(indexPath, content)
	if err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	return nil
}
