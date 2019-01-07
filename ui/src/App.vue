<template>
  <div id="app">
    <div class="ui fixed inverted menu">
      <div class="ui container">
        <a href="#" class="header item">
          <img class="logo" src="@/assets/icon.png">
          Creamy Videos
        </a>
        <a href="/" class="item">Home</a>
      </div>
    </div>

    <div class="ui main container">
      <router-view />
    </div>
  </div>
</template>

<style lang="scss">
#app>.ui.main.container {
  margin-top: 4em;
}
</style>
