<template>
  <v-card>
    <v-card-title>
      <v-col cols="2">
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn icon @click="backFolder" :loading="loadingUp" v-on="on" :disabled="back_disabled">
              <v-icon>mdi-arrow-left</v-icon>
            </v-btn>
          </template>
          <span>Back</span>
        </v-tooltip>
        <v-tooltip bottom>
          <template v-slot:activator="{ on }">
            <v-btn icon @click="fwdFolder" :loading="loadingUp" v-on="on" :disabled="fwd_disabled">
              <v-icon>mdi-arrow-right</v-icon>
            </v-btn>
          </template>
          <span>Fwd</span>
        </v-tooltip>
      </v-col>
      <v-col>
        <v-text-field
          v-model="search"
          append-icon="search"
          label="Search"
          single-line
          hide-details
        ></v-text-field>
      </v-col>
    </v-card-title>
    <v-container>
      <v-row class="mx-1 mb-1">{{ parent_folder }}</v-row>
      <v-data-table
        v-model="musicSelected"
        :headers="headers"
        :items="musicdata"
        :loading="loadingData"
        :items-per-page="itemsPerPage"
        item-key="id"
        show-select
        class="elevation-1"
        :search="search"
        :page="page"
        :pageStart="pageStart"
        :footer-props="{
          showFirstLastPage: true,
          firstIcon: 'mdi-arrow-collapse-left',
          lastIcon: 'mdi-arrow-collapse-right',
          prevIcon: 'mdi-minus',
          nextIcon: 'mdi-plus',
        }"
      >
        <template v-slot:item.actions="{ item }">
          <v-icon small class="mr-2" @click="playOrfetchSubFolder(item)">{{
            item.icon_action
          }}</v-icon>
        </template>
        <template v-slot:item.type="{ item }">
          <v-chip :color="getColorType(item.fileorfolder)" dark>{{
            item.type
          }}</v-chip>
        </template>
      </v-data-table>
    </v-container>
  </v-card>
</template>