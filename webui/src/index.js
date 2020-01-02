import './main.css';
import { Elm } from './Main.elm';
import registerServiceWorker from './registerServiceWorker';


const storageKey = "strongbox";
const flags = localStorage.getItem(storageKey);
const app = Elm.Main.init({
  flags: flags,
  node: document.getElementById('root')
});

app.ports.cache.subscribe(function(val) {
  if (val === null) {
    localStorage.removeItem(storageKey);
  } else {
    localStorage.setItem(storageKey, JSON.stringify(val));
  }
  setTimeout(function() { app.ports.onStoreChange.send(val); }, 0);
});

// Whenever localStorage changes in another tab, report it if necessary.
window.addEventListener("storage", function(event) {
  if (event.storageArea === localStorage && event.key === storageKey) {
    app.ports.onStoreChange.send(event.newValue);
  }
}, false);


registerServiceWorker();
