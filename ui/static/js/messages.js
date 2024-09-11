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
        if (change.type == "added") {
            console.log("added")
            console.log(change.doc.data());
        } else if (change.type == "modified") {
            console.log("modified")
            console.log(change.doc.data());
        } else if (change.type == "removed") {
            console.log("removed")
            console.log(change.doc.data());
        }
    });
});
