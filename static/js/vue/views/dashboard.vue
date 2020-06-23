<template>
  <v-card color="grey lighten-4" flat tile>
    <v-toolbar flat dense>
      <v-toolbar-title class="subheading grey--text">Omx Player</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="playYoutube" :loading="loadingyoutube" v-on="on">
            <v-icon>airplay</v-icon>
          </v-btn>
        </template>
        <span>Play youtube url</span>
      </v-tooltip>
      <v-tooltip bottom>
        <template v-slot:activator="{ on }">
          <v-btn icon @click="playPlaylist" :loading="loadingplaylist" v-on="on">
            <v-icon>queue_music</v-icon>
          </v-btn>
        </template>
        <span>Play Playlist</span>
      </v-tooltip>
    </v-toolbar>
    <v-col cols="12">
      <v-row>
        <v-col cols="12">
          <v-text-field v-model="uriToPlay" label="Select an URL or playlist"></v-text-field>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <v-card flat tile>
            <v-card-title>Playlist {{ListName}}</v-card-title>
            <div class="mx-4">
              <div class="subtitle-2 text--secondary">URI: {{PlayingURI}}</div>
              <div class="subtitle-2 text--secondary">Info: {{PlayingInfo}}</div>
              <div class="subtitle-2 text--secondary">Type: {{PlayingType}}</div>
              <div class="subtitle-2 text--secondary">Previous: {{PlayingPrev}}</div>
              <div class="subtitle-2 text--secondary">Next: {{PlayingNext}}</div>
            </div>
          </v-card>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12">
          <Playerbar />
        </v-col>
      </v-row>
    </v-col>
  </v-card>
</template>