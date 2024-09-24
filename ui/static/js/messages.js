import { initializeApp } from "https://www.gstatic.com/firebasejs/10.8.1/firebase-app.js";
import {
  getFirestore,
  doc,
  onSnapshot,
  query,
  collection,
  orderBy,
} from "https://www.gstatic.com/firebasejs/10.8.1/firebase-firestore.js";

const firebaseConfig = {
  apiKey: "AIzaSyC1G8vFyKjDmjF81N3E8c-IkwUeXrxuAbk",
  authDomain: "whatsapp-3a492.firebaseapp.com",
  projectId: "whatsapp-3a492",
  storageBucket: "whatsapp-3a492.appspot.com",
  messagingSenderId: "435949582458",
  appId: "1:435949582458:web:d03b812406f478401063c3",
  measurementId: "G-8ZYDT7SVDN"
};

const app = initializeApp(firebaseConfig);
const db = getFirestore(app);

const wa_id = document.querySelector("#wa_id").innerText
console.log(wa_id);
const q = query(collection(db, `contacts/${wa_id}/messages`),)
onSnapshot(q, (snapshot) => {
    snapshot.docChanges().forEach((change) => {
        const docId = change.doc.id;
        const data = change.doc.data();
        if (change.type == "added") {
            console.log("added")
            console.log(data);
            createMessage(data, docId);
            listenForStatusUpdates(docId);
        } else if (change.type == "modified") {
            console.log("modified")
            console.log(data);
        } else if (change.type == "removed") {
            console.log("removed")
            console.log(data);
        }
    });
});

function listenForStatusUpdates(messageId) {
    const statusQuery = query(collection(db, `contacts/${wa_id}/messages/${messageId}/status`));

    onSnapshot(statusQuery, (snapshot) => {
        snapshot.docChanges().forEach((change) => {
            const statusData = change.doc.data();

            if (change.type === "added" || change.type === "modified") {
                console.log(`Status update for message ${messageId}: `, statusData);
                updateMessageStatus(messageId, statusData);
            }
        });
    });
}

function getMessageTemplate() {
    let message = document.createElement('div');
    //message.id = 'message';
    message.classList.add('max-w-[70%]', 'rounded-lg', 'p-3', 'mb-2');

    let content = document.createElement('div');
    content.id = 'message_content';
    message.appendChild(content);

    let timestamp = document.createElement('p');
    timestamp.id = 'message_timestamp';
    timestamp.classList.add('text-xs', 'text-gray-500', 'text-right', 'mt-1');

    let stat = document.createElement('p');
    stat.id = 'message_status';
    stat.classList.add('text-xs', 'text-gray-500', 'text-right', 'mt-1');

    message.appendChild(timestamp);
    message.appendChild(stat);

    return message;
}

function createMessage(data, id) {
    const message = document.querySelector("#messages");
    let msgTmpl = getMessageTemplate();
    msgTmpl.id = id
    if (data.to) {
        msgTmpl.classList.add('bg-green-100', 'ml-auto');
    } else if (data.from) {
        msgTmpl.classList.add('bg-white');
    }
    let content = msgTmpl.querySelector("#message_content");
    let timestamp = msgTmpl.querySelector("#message_timestamp");
    timestamp.textContent = data.timestamp;
    switch (data.type) {
        case "text":
            let text = document.createElement('p');
            text.textContent = data.text.body;
            content.appendChild(text)
            break;
        case "image":
            let img = document.createElement('img');
            img.src = data.image.link;
            content.appendChild(img);
            break;
        default:
            console.log("new type", data.type);
            break;
    }

    message.appendChild(msgTmpl);
}

function updateMessageStatus(messageId, statusData) {
    // Find the message element in the DOM using messageId
    const messageElement = document.getElementById(messageId);

    if (messageElement) {
        const statusElement = messageElement.querySelector("#message_status");
        if (statusElement) {
            // Update the status element with the new status (e.g., "sent", "delivered", "failed")
            statusElement.textContent = statusData.status;

            // Optionally, add some visual indicators based on status
            //if (statusData.status === "delivered") {
                //messageElement.classList.add("delivered");
            //} else if (statusData.status === "failed") {
                //messageElement.classList.add("failed");
            //}
        }
    }
}
