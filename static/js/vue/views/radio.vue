<template>
  <v-container pa-1>
    <v-skeleton-loader
      :loading="radioloading"
      :transition="transition"
      height="94"
      type="list-item-three-line"
    >
      <v-card color="grey lighten-4" flat tile>
        <v-toolbar flat dense>
          <v-toolbar-title class="subheading grey--text">Radio</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-tooltip bottom>
            <template v-slot:activator="{ on }">
              <v-btn icon @click="prepareInsert" v-on="on">
                <v-icon>mdi-plus</v-icon>
              </v-btn>
            </template>
            <span>New Radio</span>
          </v-tooltip>
        </v-toolbar>
        <v-container>
          <v-list dense nav>
            <template v-for="plitem in radio">
              <v-list-item :key="plitem.id">
                <v-list-item-content>
                  <v-list-item-title>{{ plitem.title }}</v-list-item-title>
                  <v-list-item-title>{{
                    plitem.description
                  }}</v-list-item-title>
                  <v-list-item-title>{{ plitem.genre }}</v-list-item-title>
                  <v-list-item-title>{{ plitem.uri }}</v-list-item-title>
                  <v-row>
                    <v-btn
                      icon
                      text
                      :key="plitem.id"
                      @click="askForPlayItem(plitem)"
                      ><v-icon>library_music</v-icon>
                    </v-btn>
                    <v-spacer></v-spacer>
                    <v-btn icon text @click="prepareEdit(plitem)"
                      ><v-icon>mdi-circle-edit-outline</v-icon>
                    </v-btn>
                    <v-btn icon text @click="askForDeleteItem(plitem)"
                      ><v-icon>mdi-delete-forever-outline</v-icon>
                    </v-btn>
                  </v-row>
                </v-list-item-content>
              </v-list-item>
            </template>
          </v-list>
          <v-divider></v-divider>
          <v-row justify="center">
            <v-btn icon text @click="loadMore" :disabled="last_radio_fetch"
              >More<v-icon>more_horiz</v-icon>
            </v-btn>
          </v-row>
        </v-container>
      </v-card>
    </v-skeleton-loader>
    <v-container>
      <v-dialog v-model="dialogItemSelected" persistent max-width="290">
        <v-card>
          <v-card-title class="headline">Question</v-card-title>
          <v-card-text
            >Do you want to {{ selected_item.action_name }} the radio "{{
              selected_item.itemquestion
            }}"?</v-card-text
          >
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="continueSelectedItem"
              >OK</v-btn
            >
            <v-btn
              color="green darken-1"
              text
              @click="dialogItemSelected = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
      <v-dialog v-model="dialogInsertEdit" persistent max-width="290">
        <v-card>
          <v-container>
            <v-col cols="12">
              <v-row justify="space-around">
                <v-card-title class="headline">{{dialog_title}}</v-card-title>
                <v-text-field
                  label="Name"
                  v-model="selected_item.title"
                  :rules="rules.name"
                  required
                ></v-text-field>
                <v-text-field
                  label="URI"
                  v-model="selected_item.uri"
                  :rules="rules.URI"
                  required
                ></v-text-field>
                <v-text-field
                  label="Description"
                  v-model="selected_item.description"
                ></v-text-field>
              </v-row>
            </v-col>
          </v-container>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="green darken-1" text @click="continueSelectedItem"
              >OK</v-btn
            >
            <v-btn color="green darken-1" text @click="dialogInsertEdit = false"
              >Cancel</v-btn
            >
          </v-card-actions>
        </v-card>
      </v-dialog>
    </v-container>
  </v-container>
</template>