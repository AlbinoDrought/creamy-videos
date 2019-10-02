<template>
  <div id="app">
    <div class="ui fixed inverted main menu">
      <div class="ui container">
        <router-link :to="{ name: 'home' }" class="header item">
          <img class="logo" src="@/assets/icon.png">
          Creamy Videos
        </router-link>
        <router-link :to="{ name: 'home' }" class="item">
          Home
        </router-link>
        <router-link v-if="!readOnly" :to="{ name: 'upload' }" class="item">
          Upload
        </router-link>

        <!-- this search menu is hidden in mobile view -->
        <div class="not-small right menu">
          <div class="borderless item" v-if="$route.meta.sortable">
            <sort-dropdown v-model="sortKey" />
          </div>

          <div class="borderless item">
            <div class="ui inverted transparent icon input" @keypress.enter="search">
              <!-- not using v-model because of https://github.com/vuejs/vue/issues/8231 -->
              <!-- `searchText` was not properly updating when pressing enter on mobile -->
              <input
                aria-label="Search"
                type="text"
                placeholder="Search..."
                :value="searchText"
                @input="searchText = $event.target.value"
                @keypress.enter="search"
              >
              <i class="search link icon" @click.prevent="search" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- separate search menu for mobile -->
    <!-- design inspired from one of my favorite sites ;) -->
    <div class="ui only-small fluid fixed inverted menu search">
      <div class="borderless item">
        <div class="ui inverted transparent icon input">
          <!-- not using v-model because of https://github.com/vuejs/vue/issues/8231 -->
          <!-- `searchText` was not properly updating when pressing enter on mobile -->
          <input
            type="text"
            placeholder="Search..."
            :value="searchText"
            @input="searchText = $event.target.value"
            @keypress.enter="search"
          >
          <i class="search link icon" @click.prevent="search" />
        </div>
      </div>
    </div>

    <div class="ui main container">
      <router-view :key="searchKey" />
    </div>
  </div>
</template>

<script>
import SortDropdown from '@/components/SortDropdown.vue';
import sortOptions from '@/sortOptions';

export default {
  components: {
    SortDropdown,
  },

  computed: {
    readOnly() {
      return this.$store.getters.readOnly;
    },

    sortKey: {
      get() {
        return this.$route.query.sort || sortOptions[0].key;
      },

      set(sortKey) {
        this.$router.push({
          name: this.$route.name,
          query: {
            ...this.$route.query,
            sort: sortKey,
          },
        });
      },
    },
  },

  data() {
    return {
      searchText: '',
      searchKey: 0,
    };
  },

  metaInfo: {
    titleTemplate: titleChunk => (titleChunk ? `${titleChunk} | creamy-videos` : 'creamy-videos'),
  },

  methods: {
    search() {
      this.$router.push({
        name: 'search',
        query: {
          text: this.searchText,
          sort: this.sortKey,
        },
      });
      // increment the searchKey so it is unique,
      // and anything that uses it gets remounted.
      // this is a semi-hack to force multiple searches
      // of the same search term to actually search again.
      this.searchKey += 1;
    },
  },

  mounted() {
    if (this.$route.query.text) {
      this.searchText = this.$route.query.text;
    }
  },
};
</script>

<style lang="scss">
// grabbed during runtime
$main-menu-height: 53px;
$mobile-search-menu-height: 40px;
$base-top-margin: $main-menu-height + 6px;
$mobile-top-margin: $base-top-margin + $mobile-search-menu-height;

#app .main.menu {
  height: $main-menu-height;

  .header.item img {
    margin-right: 10px;
  }
}

#app>.ui.main.container {
  // the fixed semantic-ui menu doesn't
  // pad the contents below it.
  // without this margin-top, some contents
  // would be obscured
  margin-top: $base-top-margin;
}

#app .only-small {
  display: none;
}

#app .search.menu {
  top: $main-menu-height;
  height: $mobile-search-menu-height;
  border-bottom: 1px solid rgba(255, 255, 255, 0.15);
}

@media only screen and (max-width: 767px) {
  #app .main.menu .item::before {
    // disable semantic-ui "pseudo-borders" in mobile
    background: none;
  }

  #app .main.menu .header.item {
    // align menu items to the right in mobile
    flex-grow: 1;
  }

  #app>.ui.main.container {
    // when we show the mobile search menu,
    // offset main container contents even more
    // so nothing is obscured
    margin-top: $mobile-top-margin;
  }

  #app .only-small {
    display: block;
  }

  #app .not-small {
    display: none;
  }
}
</style>
