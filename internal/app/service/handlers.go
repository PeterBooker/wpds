package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

func appHandler(w http.ResponseWriter, r *http.Request) {

	render(w, "app", nil)

}

func apiPageSearchViewHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	searchID := chi.URLParam(r, "id")

	w.Write([]byte(`{
			"title": "Search View - WPDS",
			"heading": "Search View",
			"summary": {
				"id": ` + searchID + `,
				"directory": "Plugins",
				"term": "^$//fake-.regex%/",
				"top": [
				{
					"matches": 5,
					"slug": "rest-api",
					"installs": "40,000+"
				},
				{
					"matches": 9,
					"slug": "wptoandroid",
					"installs": "30+"
				},
				{
					"matches": 11,
					"slug": "custom-contact-form",
					"installs": "60,000+"
				},
				{
					"matches": 2,
					"slug": "ninja-forms",
					"installs": "50+"
				},
				{
					"matches": 2,
					"slug": "simple-customizer",
					"installs": "500+"
				}
				] 
			},
			"matches": [
				{
					"file": "appmaker-woocommerce-mobile-app-manager/lib/appmaker-wp/endpoints/appmaker/class-appmaker-wc-rest-backend-posts-controller.php",
					"line": 705,
					"content": "$date_data = rest_get_date_with_gmt( $request['date'] );"
				},
				{
					"file": "appmaker-woocommerce-mobile-app-manager/lib/appmaker-wp/endpoints/appmaker/class-appmaker-wc-rest-backend-posts-controller.php",
					"line": 705,
					"content": "$date_data = rest_get_date_with_gmt( $request['date'] );"
				},
				{
					"file": "appmaker-woocommerce-mobile-app-manager/lib/appmaker-wp/endpoints/appmaker/class-appmaker-wc-rest-backend-posts-controller.php",
					"line": 705,
					"content": "$date_data = rest_get_date_with_gmt( $request['date'] );"
				},
				{
					"file": "appmaker-woocommerce-mobile-app-manager/lib/appmaker-wp/endpoints/appmaker/class-appmaker-wc-rest-backend-posts-controller.php",
					"line": 705,
					"content": "$date_data = rest_get_date_with_gmt( $request['date'] );"
				},
				{
					"file": "appmaker-woocommerce-mobile-app-manager/lib/appmaker-wp/endpoints/appmaker/class-appmaker-wc-rest-backend-posts-controller.php",
					"line": 705,
					"content": "$date_data = rest_get_date_with_gmt( $request['date'] );"
				},
				{
					"file": "appmaker-woocommerce-mobile-app-manager/lib/appmaker-wp/endpoints/appmaker/class-appmaker-wc-rest-backend-posts-controller.php",
					"line": 705,
					"content": "$date_data = rest_get_date_with_gmt( $request['date'] );"
				},
				{
					"file": "appmaker-woocommerce-mobile-app-manager/lib/appmaker-wp/endpoints/appmaker/class-appmaker-wc-rest-backend-posts-controller.php",
					"line": 705,
					"content": "$date_data = rest_get_date_with_gmt( $request['date'] );"
				},
				{
					"file": "appmaker-woocommerce-mobile-app-manager/lib/appmaker-wp/endpoints/appmaker/class-appmaker-wc-rest-backend-posts-controller.php",
					"line": 705,
					"content": "$date_data = rest_get_date_with_gmt( $request['date'] );"
				},
				{
					"file": "appmaker-woocommerce-mobile-app-manager/lib/appmaker-wp/endpoints/appmaker/class-appmaker-wc-rest-backend-posts-controller.php",
					"line": 705,
					"content": "$date_data = rest_get_date_with_gmt( $request['date'] );"
				}
			]
		}
		`))

}

func apiPageSearchListHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write([]byte(`{
			"title": "Search List - WPDS",
			"heading": "Search List",
			"latest_searches": [
				{
					"id": "5",
					"directory": "plugins",
					"completed": "2017/10/13 11:49",
					"term": "trash_post",
					"num_matches": 7654
				},
				{
					"id": "4",
					"directory": "themes",
					"completed": "2017/10/13 10:12",
					"term": "$//+.more/regex?.+/",
					"num_matches": 217
				},
				{
					"id": "3",
					"directory": "plugins",
					"completed": "2017/10/13 09:36",
					"term": "^$//fake-.regex%/",
					"num_matches": 786
				},
				{
					"id": "2",
					"directory": "themes",
					"completed": "2017/10/13 09:05",
					"term": "prefix_func()",
					"num_matches": 27
				},
				{
					"id": "1",
					"directory": "plugins",
					"completed": "2017/10/13 07:53",
					"term": "prefix_func()",
					"num_matches": 13
				}
			],
			"recurring_searches": [
				{
					"title": "Test",
					"content": "<div class=\"loading-container\"><div class=\"loading\"></div></div>"
				}
			]
			}
		`))

}

func apiPageDashboardHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write([]byte(`{"page":
		"title": "WPDS - WordPress Directory Slurper",
		"heading": "Dashboard",
		"boxes": [
			{
				"title": "Test",
				"content": "<div class="loading-container"><div class="loading"></div></div>"
			},
			{
				"title": "Test",
				"content": "<div class="loading-container"><div class="loading"></div></div>"
			}
		]
		}
	`))

}

// Render Web Page HTML
func render(w http.ResponseWriter, template string, data interface{}) {

	err := t.ExecuteTemplate(w, template, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
