<template>
  <div id="app">
    <div class="ui fixed inverted menu">
      <div class="ui container">
        <a href="#" class="header item">
          <img class="logo" src="@/assets/icon.png">
          Creamy Videos
        </a>
        <router-link :to="{ name: 'home' }" class="item">
          Home
        </router-link>
        <router-link :to="{ name: 'upload' }" class="item">
          Upload
        </router-link>

        <div class="right menu">
          <div class="borderless item">
            <div class="ui inverted transparent icon input" @keypress.enter="search">
              <input type="text" placeholder="Search..." v-model="searchText">
              <i class="search link icon" @click.prevent="search" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="ui main container">
      <router-view />
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      searchText: '',
    };
  },
  
  methods: {
    search() {
      this.$router.push({
        name: 'search',
        query: {
          text: this.searchText,
        },
      });
    },
  },
};
</script>

<style lang="scss">
#app>.ui.main.container {
  margin-top: 4em;
}
</style>
