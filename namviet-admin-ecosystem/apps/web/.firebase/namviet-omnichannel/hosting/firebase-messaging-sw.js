importScripts('https://www.gstatic.com/firebasejs/10.8.0/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/10.8.0/firebase-messaging-compat.js');

// Configuration from your .env.local (manually hardcoded for SW or fetched, but hardcoded is simpler for Service Worker)
// In production, you might want to dynamically inject these variables
firebase.initializeApp({
  apiKey: "AIzaSyAx3E4bOyZ24rJ0JhoLNAH-Ih_kuYJvacg",
  authDomain: "namviet-omnichannel.firebaseapp.com",
  projectId: "namviet-omnichannel",
  storageBucket: "namviet-omnichannel.firebasestorage.app",
  messagingSenderId: "975139872486",
  appId: "1:975139872486:web:715d3bf2f81743f4d32527",
});

const messaging = firebase.messaging();

messaging.onBackgroundMessage((payload) => {
  console.log('[firebase-messaging-sw.js] Received background message ', payload);
  
  const notificationTitle = payload.notification.title;
  const notificationOptions = {
    body: payload.notification.body,
    icon: '/logo.png',
  };

  self.registration.showNotification(notificationTitle, notificationOptions);
});
