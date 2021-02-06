package controllers

import "github.com/victorsteven/fullstack/api/middlewares"

func (s *Server) initializeRoutes() {

	//TESTING API EXT GATEWAY
	s.Router.HandleFunc("/wa/send", middlewares.SetMiddlewareJSON(s.SendWhatsapp)).Methods("POST")

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/authorized", middlewares.SetMiddlewareAuthentication(s.CheckAuth)).Methods("GET")
	s.Router.HandleFunc("/phone/otp", middlewares.SetMiddlewareJSON(s.PhoneOTPRequest)).Methods("POST")
	s.Router.HandleFunc("/login/phone", middlewares.SetMiddlewareJSON(s.LoginPhone)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	s.Router.HandleFunc("/home", middlewares.SetMiddlewareJSON(s.GetHome)).Methods("GET", "OPTIONS")

	//product
	// s.Router.HandleFunc("/product/recomedation", middlewares.SetMiddlewareJSON(s.GetProductRecomendation)).Methods("GET")
	s.Router.HandleFunc("/product/{slug}", middlewares.SetMiddlewareJSON(s.GetProduct)).Methods("GET")
	s.Router.HandleFunc("/search", middlewares.SetMiddlewareJSON(s.GetProducts)).Methods("POST")

	//cart
	s.Router.HandleFunc("/address", middlewares.SetMiddlewareAuthentication(s.GetAddress)).Methods("GET")
	s.Router.HandleFunc("/address", middlewares.SetMiddlewareAuthentication(s.AddAddress)).Methods("POST")
	s.Router.HandleFunc("/address/update", middlewares.SetMiddlewareAuthentication(s.UpdateAddress)).Methods("POST")
	s.Router.HandleFunc("/address", middlewares.SetMiddlewareAuthentication(s.DelAddress)).Methods("DELETE")
	s.Router.HandleFunc("/region", middlewares.SetMiddlewareJSON(s.GetRegion)).Methods("GET")
	s.Router.HandleFunc("/cart", middlewares.SetMiddlewareJSON(s.AddCart)).Methods("POST")
	s.Router.HandleFunc("/couriers", middlewares.SetMiddlewareJSON(s.GetCourier)).Methods("POST")

	//order
	s.Router.HandleFunc("/checkout", middlewares.SetMiddlewareJSON(s.Checkout)).Methods("POST")
	s.Router.HandleFunc("/orders", middlewares.SetMiddlewareAuthentication(s.GetOrders)).Methods("GET")
	s.Router.HandleFunc("/order/{orderID}", middlewares.SetMiddlewareJSON(s.OrderDetail)).Methods("GET")

	s.Router.HandleFunc("/payments", middlewares.SetMiddlewareJSON(s.GetPayments)).Methods("GET")
}
