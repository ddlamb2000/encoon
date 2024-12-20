import { building } from "$app/environment";
import { connect } from "$lib/kafka";

if (!building) {
    connect().then(() => {
        console.log("Connected to Kafka.");
    }).catch((error) => {
        console.error("There's an error connecting to Kafka:", error);
    });
}
