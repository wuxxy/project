use axum::response::IntoResponse;
use axum::routing::{get, post};
use serde_json::json;
use socketioxide::extract::SocketRef;
use socketioxide::SocketIo;
struct Message {
    event: String,
    data: String
}
#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let (layer, io) = SocketIo::new_layer();

    // Register a handler for the default namespace
    io.ns("/", |s: SocketRef| {
        s.on("message", |s: SocketRef| {
            println!("{}", s.id);
            s.emit("message-back", json!({
                "echo": "test"
            }).to_string().as_str()).ok();
        });
    });
    let app = axum::Router::new()
        .route("/", get(|| async { "Hello, World!" }))
        .layer(layer)
        .route("/send_message", post({
            let io = io.clone();
            move || {
                let io = io.clone();
                async move {
                    io.broadcast()
                        .emit("message", "weirdo")
                        .await
                        .expect("Failed to emit message");
                    "Message sent".into_response()
                }
            }
        }));
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();

    Ok(())
}
