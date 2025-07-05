package http

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "github.com/aventhis/go-order-service/internal/service"
)

type Handler struct {
    orderService service.OrderServiceInterface
}

func NewHandler(orderService service.OrderServiceInterface) *Handler {
    return &Handler{
        orderService: orderService,
    }
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    orderUID := vars["id"]

    order, err := h.orderService.GetByID(orderUID)
    if err != nil {
        http.Error(w, "Order not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(order)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
    html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Order Info</title>
    </head>
    <body>
        <h1>Enter Order ID</h1>
        <input type="text" id="orderID" placeholder="Order ID">
        <button onclick="getOrder()">Search</button>
        <div id="result"></div>

        <script>
        function getOrder() {
            const id = document.getElementById('orderID').value;
            fetch('/order/' + id)
                .then(response => response.json())
                .then(data => {
                    document.getElementById('result').innerText = 
                        JSON.stringify(data, null, 2);
                })
                .catch(() => {
                    document.getElementById('result').innerText = 
                        'Order not found';
                });
        }
        </script>
    </body>
    </html>
    `
    w.Header().Set("Content-Type", "text/html")
    w.Write([]byte(html))
}

func (h *Handler) InitRoutes() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/", h.Index).Methods("GET")
    r.HandleFunc("/order/{id}", h.GetOrder).Methods("GET")
    return r
}