use actix_web::{post, web, App, HttpResponse, HttpServer, Responder, get};
use serde::{Deserialize, Serialize}; // Import both Serialize and Deserialize

#[derive(Debug, Serialize, Deserialize)] // Derive both Serialize and Deserialize
struct WeatherData {
    description: String,
    country: String,
    weather: String,
}

#[get("/health")]
async fn health_check() -> impl Responder {
    HttpResponse::Ok().body("OK")
}


#[post("/input")]
async fn handle_input(info: web::Json<WeatherData>) -> impl Responder {
    let client = reqwest::Client::new();
    println!("📦 Enviando datos al API REST: {:?}", info);
    let res = client.post("http://go-api-rest-service.grpc-namespace.svc.cluster.local:8081/input")
        .json(&*info)
        .send()
        .await;

    match res {
        Ok(_) => HttpResponse::Ok().json(info.0),
        Err(e) => HttpResponse::InternalServerError().body(format!("Error: {}", e)),
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("🚀 Iniciando servidor en 0.0.0.0:8080"); // <-- Este log

    HttpServer::new(|| {
        App::new()
            .service(handle_input)
            .service(health_check)
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}