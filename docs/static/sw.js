// Self-destroying service worker.
//
// The previous revive.run site (Gatsby) registered a service worker at /sw.js
// on visitors' browsers. Serving this file at the same scope makes those
// browsers pick it up as an update: on install it takes over immediately, on
// activate it unregisters itself and reloads every open tab so pages are
// served from the network again instead of the stale Gatsby cache.
//
// Standard pattern: https://github.com/NekR/self-destroying-sw
self.addEventListener('install', () => {
  self.skipWaiting();
});

self.addEventListener('activate', () => {
  self.registration
    .unregister()
    .then(() => self.clients.matchAll({ type: 'window' }))
    .then((clients) => {
      clients.forEach((client) => {
        if (client.url && 'navigate' in client) {
          client.navigate(client.url);
        }
      });
    });
});
