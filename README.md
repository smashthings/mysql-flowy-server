# MySQL Backed Flowy Server

Flowy is a workflowy clone by @suyash . I'm a big fan however the backend servers in their below repo didn't build or used GCP. I'm mostly on-prem so this repo rebuilds the backend from the appengine server into a MySQL backed DB.

## Notes

- As this repo acts as a server and specifically built for an on-prem trusted network, *there's no authorisation or authentication*
- The concept of API Key has been altered to instead act as a database table key so that multiple users can use this server as a backend 
- There's no further security or other elements involved. This isn't getting scanned or whatever
- Database entries are base64 encoded for handling text data, hence part of the DB is not going to be human readable (text entries)

## Development

Use `docker-compose up` to run all the services needed for local development. You'll need to restart services to rerun any compiled binaries due to dockerisms

## Credits & References

suyash - the author of Flowy and the servers that this repo is based on => https://github.com/suyash  
Flowy - the frontend that runs as part of the server => https://github.com/suyash/flowy  
Flowy Servers - the original still available backend servers for flowy that this repo is based on => https://github.com/suyash/flowy-servers  

