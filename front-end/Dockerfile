# Usa l'immagine ufficiale di NGINX
FROM nginx:alpine

# Rimuove i file default (opzionale ma pulito)
RUN rm -rf /usr/share/nginx/html/*

# Copia i file del tuo front-end nella cartella pubblica di nginx
COPY ./ /usr/share/nginx/html/

# Espone la porta 80 (HTTP)
EXPOSE 80

# Il comando CMD è già definito nell'immagine di base:
# CMD ["nginx", "-g", "daemon off;"]
