# live-omxctrl
Programma per gestire files e urls lanciati sul Raspi Pi4 collegato allo stereo.
Lo scopo è quello di accedere il Pi4, lo stereo e di selezionare una risorsa, che sia una url un video
o un file musicale, attraverso una interfaccia web che funzioni anche, e sopratutto, su un iphone.

## Uso abituale
- Lanciare una nuova uri del tubo.
- Lanciare una uri ascoltata ieri.
- Lanciare un video memorizzato sul pi4

## comandi frequenti
Per stoppare il sevice si usa:

    sudo systemctl stop live-omxctrl

Vue si usa nella root del progetto con:

    /usr/local/bin/vuetojs.bin -vue ./static/js/vue/views/dashboard.vue

## Sviluppo di live-omxctrl
Su winows 10 uso visual code in modalità remota con ssh:pi4.
Durante lo sviluppo uso la porta 5549, mentre quella del service è la 5548, così posso ascoltare
la musica durante la programmazione.
Da notare che quasta modalità funziona egregiamente su un Raspy 4, ma sun Raspi 3 fa molta più fatica.
Quindi il Raspy 3 lo uso solo come deployment target, questo solo in teoria in quanto non
ho intenzione di fare il deployment live-omxctrl su pi3, ma uso un altro service live-streamer
su questo target.

## Deployment su arm direttamente
In un colpo: ./publish-pi4.sh

In dettaglio quello che viene eseguito:
go build -o live-omxctrl.bin
cd deploy
./deploy.bin -target pi4 -outdir ~/app/live-omxctrl/zips/
cd ~/app/live-omxctrl/
./update-service.sh

## Sviluppo su pi4 arm
Apri vscode nella directory remota (nota come la costruzione della directory. /go/ è la differenza tra sviluppo e deploy):
/home/igor/app/go/live-omxctrl/


### Deployment dettagli e preparazione
- Sul laptop occorre il file copy_app_to_pi4.sh posizionato nella dir ../deployed
- Su Pi4 occorre la directory /home/igors/app/live-omxctrl con all` interno il file update-service.sh
così come la dir /home/igors/app/live-omxctrl/zips

## Preparazione target pi4
Sul firewall di pi4 ho aperto la porta 5548 in intranet con:

    ufw allow from 192.168.2.0/24 to any port 5548

Ora va installata la app. Uso la dir:
~/app/go$mkdir live-omxctrl\zips
copio lo zip deployed locale in live-omxctrl\zips con ./copy_app_to_pi4.sh
copio ./update-service.sh in ~/app/go/live-omxctrl e lo lancio per scompattare lo zip nella dir ./current
Poi si va ./current e si prova il service con: ./live-omxctrl.bin

Poi si mette il programma live-omxctrl.bin come service di sistema.

    sudo nano /lib/systemd/system/live-omxctrl.service

Abilitare il service:
sudo systemctl enable live-omxctrl.service
Ora si fa partire il service (resistente al reboot):

    sudo systemctl start live-omxctrl

Per vedere i logs si usa:

    sudo journalctl -f -u live-omxctrl

## Sviluppo su Pi4 
È possibile sviluppare il software direttamente su pi4 usando l'extension remote ssh
Basta installare su pi4 golang arm6. Il vantaggio, oltre al deployment ancora più semplice,
dovrebbe essere la possibilità di programmare direttamente dbus senza usare la shell.
Dbus su windows è molto diverso.
Per l'extension vuetojs non funziona su remoto in quanto ho impacchettato un tool in formato
windows. L'ho compilato separatamente e messo in /usr/local/bin/vuetojs.bin 
Per usarlo si scrive nella bash:

    /usr/local/bin/vuetojs.bin -vue ./static/js/vue/views/dashboard.vue

### Dbus
Per comunicare con il processo omxplayer si usa il dbus con il suo protocollo (solo linux, pi4).
Esiste già lo script bash che esegue i comandi, file bash/dbuscontrol.sh.
Esso setta 2 variabili di ambiente:
DBUS_SESSION_BUS_ADDRESS
DBUS_SESSION_BUS_PID
Poi lancia il comando specifico, per esempio per la duration:
dbus-send --print-reply=literal --session --reply-timeout=500 --dest=org.mpris.MediaPlayer2.omxplayer /org/mpris/MediaPlayer2 org.freedesktop.DBus.Properties.Get string:"org.mpris.MediaPlayer2.Player" string:"Duration"
Questa serie di comandi dovrebbe essere possibile in golang.
Ci sono due files legati all'utente:

/tmp/omxplayerdbus.igors => esempio unix:abstract=/tmp/dbus-1OTLRLIFgE,guid=39be549b2196c379ccdf29585ed9674d
Qui dovrebbe descritto quale socket ascolta il server, un socket del tipo abstract che sembra
un path del file di sistema ma non lo è

/tmp/omxplayerdbus.igors.pid  => esempio 1379
Questo dovrebbe essere il pid del service daemon dbus-daemon:

    1 S  1001  1379     1  0  80   0 -  1604 do_epo ?        00:00:00 dbus-daemon

I due files /tmp/omxplayerdbus.igors vengono creati quando si fa partire il player.

## Service Config
Questo il conetnuto del file che compare con:
sudo nano /lib/systemd/system/live-omxctrl.service
Poi si fa l'enable:
sudo systemctl enable live-omxctrl.service
E infine lo start:
sudo systemctl start live-omxctrl
Logs sono disponibili con:
sudo journalctl -f -u live-omxctrl

Qui segue il contenuto del file live-omxctrl.service
Nota il Type=idle che è meglio di simple in quanto così 
viene fatto partire quando anche la wlan ha ottenuto l'IP intranet
per consentire l'accesso.

Il service lo faccio andare con l'utente Pi in quando il dbus ha un'istanza
sola ed è legata all'utente. Il db deve appartenere a pi.
Se sviluppo con igors usando un'altra porta,
dbus usa l'istanza legata all'utente. Questo vuol dire che il service
viene lanciato sotto l'utente pi, mentre lo sviluppo si svolge sotto igors.
Così non si hanno quasi mai conflitti. Di rado mi va a cambiare le istanze del service
anche se va sotto un altro utente, e da qui il dbus proprio non mi convince neanche un po'.
Con una porta del processo sarebbe tutto molto più comodo e separato per avere istanze
multiple. 

-------------------------------- file content
```
[Install]
WantedBy=multi-user.target

[Unit]
Description=live-omxctrl service
ConditionPathExists=/home/igors/app/live-omxctrl/current/live-omxctrl.bin
After=network.target

[Service]
Type=idle
User=pi
Group=pi
LimitNOFILE=1024

Restart=on-failure
RestartSec=10
startLimitIntervalSec=60

WorkingDirectory=/home/igors/app/live-omxctrl/current/
ExecStart=/home/igors/app/live-omxctrl/current/live-omxctrl.bin

# make sure log directory exists and owned by syslog
PermissionsStartOnly=true
ExecStartPre=/bin/mkdir -p /var/log/live-omxctrl
ExecStartPre=/bin/chown pi:pi /var/log/live-omxctrl
ExecStartPre=/bin/chmod 755 /var/log/live-omxctrl
StandardOutput=syslog
StandardError=syslog

```
------------------------------------------- end file content

## Shutdown e Reboot
Per avere i comandi di reboot e shutdown sull'interfaccia web bisogna usare sudo senza password.
Nel terminale di pi4 Si lancia: 
sudo visudo
e qui si mette la linea:

    %adm  ALL=(ALL) NOPASSWD:ALL

che vuol dire che se l'utente è del gruppo adm chiama sudo senza password.

## Player
Pe fare andare una play list penso che ci voglia un blocking del 
comando che gira in una goroutine. Quando poi finisce  fa partire
la prossima. Altrimenti è difficile sapere quando il player finisce.
Senza questo non si sa quando deve ricominciare. NextTitle con il dbus
funziona solo se il player è attivo. Questo non è più vero quando la track è finita.
Penso anche che ci voglia un websocket per lo stato. Quando una track finisce
il player è off, ma il browser rimane in stato green.

### urls usate durante lo sviluppo
    //u = "http://stream.srg-ssr.ch/m/rsc_de/aacp_96"
    //u = "/home/igors/Music/tubo/Gianna Nannini - Fenomenale (Official Video)-HKwWcJCtwck.mp3"
    //u = "https://www.tubo.com/watch?v=3czUk1MmmvA"
    //u = "`tubofff-dl -f mp4 -g https://www.tubo.com/watch?v=3czUk1MmmvA`"

## Sqlite
Sqlite è il database dove vengono salvati tutti i dati.
Per vedere come si usa sqlite in full search mode vedi
https://github.com/aaaasmile/iol-importer/blob/master/Readme_iol-vienna.txt
Su raspberry il database si può gestire con interfaccia grafica usando sqlitebrowser.
sudo apt-get install sqlite3
sudo apt-get install sqlitebrowser
Per fare andare sqlitebrowser bisogna far partire Xming server in windwos.
Poi in WLC si lancia:
export DISPLAY=localhost:0.0
ssh -Y pi4
Il db del service deve avere come owner pi, in quanto uso l'utente pi per il deployment.
Ho avuto dei problemi con il database in quanto il service era sempre readonly (utente pi).
Quello che mi ha risolto il problema è stato il comando:
igors@pi4:~/app/go $ chmod 777 db
Vale a dire la directory dove si trova il file sqlite deve essere accessibile.

Ho cambiato la struttura del database usando una copia su windows. L'ho rispedito indietro con:
rsync -av ./test-data.db  igors@pi4:/home/igors/projects/go/live-omxctrl/db/test/test-data.db
Per la produzione
rsync -av ./test-data.db  igors@pi4:/home/igors/app/db/service-liveomx.db

## Soundcloud
Iniziato in qualche modo, ma non finita.
Si prende come riferimento l'extension di mopidy: 
https://github.com/mopidy/mopidy-soundcloud
Clientid non si ottiene da sound cloud. Vedi il file soundcloud_test.go
per vedere com si usa la API.
Si lancia come il video dei tube basta avere una trackid
Per esempio: https://api.soundcloud.com/tracks/62576046

## Terminare il player
Una funzionalità che ho fatto molta fatica a capire, in quanto non ha mai funzionato a dovere.
Se uso il dbus, bisogna trovare il file corretto per mandare il segnale 15. 
Meglio usare il process kill del comando exec. Il problema è che allo start vengono
creati due processi e vanno entrambi killati. Il Kill sul parent non è sufficiente.
Si usa:
cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
Nella funzione execCommand() c'è un pattern che usa cmd Start e Wait. Wait blocca
fino alla fine dell'esecuzione. Per questo è messa in una go sub routine e setta il channel done.
Il channell done è buffered in quanto se arriva prima un segnale di kill, quando Wait()
finisce il channel done sarebbe settato senza essere consumato e genera un leak.

## Windows
Si potrebbe far funzionare il service anche sotto windows, per esempio installando VLC
e usare il seguente comando per acoltare una radio:
"C:\Program Files\VideoLAN\VLC\vlc.exe" -I dummy --dummy-quiet http://stream.srg-ssr.ch/m/rsc_de/mp3_128
In questo caso VLC funziona in background e bisogna trovare il modo di controllarlo.
Per il controllo da remoto si può usare il socket.
Interfaccia "-I rc --rc-host=localhost:9876".
C'è un gist in c# che può essere usato come riferimento:
https://gist.github.com/SamSaffron/101357


## TODO
- Music: Nella navigazione del musica sui files c'è bisogno di uno stack per tornare indietro o sopra
- Music: si può far partire un file, bisgnerebbe anche stopparlo nella stessa view
- Music: info meta del file non vengono lette
- Music: seleziona dei files o folder da aggiungere alla Playlist current
- Music: mostra le playlists e lancia la playlist

- Mute e unmute mi cancella le info come title e description. Questo perché lo stato
viene riscritto completamente, invece di essere incrementale. Setstate va chiamata solo all'interno
di listenStateAction.
- Da rividere anche lo sleep del check status.
- Mettere un icona che indica lo stato di connession col ws socket
- Il collegamento con ws socket va fatto non solo al reload, ma anche quando il server si disconnete eil client manda un comando.
- Posizione e durata
- Play del soundcloud
- files audio memorizzati su pi4 da mettere nel db e in una view da cercare
- Favoriti
- Play della playlist
- Previous, Random e riciclo

- Supporto per lo stream della uri [NR - uso il progetto stream apposito]
- Play radio e podcast [DONE]
- lista di video [DONE]
- dbus è da rivedere. Se c'è gia un player in funzione, la seconda istanza su una porta diversa (dev)
va a cambiare l'istanza del service. Andrebbe magari anche isolata.[DONE]
- Le icone di mute non vanno bene. Ne serve una che toggle tra volume_down e volume_off
Le icone del volume indicano lo stato attuale e non l'azione. [DONE]
Invece il volume si fa con una slidebar. [REJ per il momento no]
- Il mutex sullo state del player dovrebe essere superfluo. Non lo è in quanto ho multipli device
su un singolo player. [DONE]
- Non mi piace la struttura di OmxPlayer è troppo grande e ridistribuita in strutture
più piccole e semplici. Per eesempio execCommand non può usare i parametri
della struttura OmxPlayer, che è un singleton, ma va usato lo stack con i parametri sulla funzione.
Il track di comandi multipli va fatto così come il kill di tutti comandi fatti partire. [DONE]
- Quando parte con una URI tipo tubo, ci deve essere associata una playlist con
un titolo solo, altrimenti mi parte un altro omx [DONE]
- La url può essere un mp3 ma anche un video. Quindi non si lancia tubo in 
automatico. [DONE]
- Perde lo stato Title e Description dopo mute/unmute [DONE]
- restart del service via os view [DONE]
- Database con la history delle url usate [DONE]

Prova a fare andare qualcosa di simile a:
/home/igors/Music/tubo/gianna-fenomenale.mp3

Per stoppare il sevice si usa:
sudo systemctl stop live-omxctrl

Per generare js
 /usr/local/bin/vuetojs.bin -vue ./static/js/vue/views/dashboard.vue 


