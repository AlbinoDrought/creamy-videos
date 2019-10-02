<template>
  <select v-model="sortOption" :class="{ fluid }">
    <option
      v-for="option in sortOptions"
      :key="option.key"
      :value="option"
    >Sort: {{ option.name }}</option>
  </select>
</template>

<script>
import sortOptions from '@/sortOptions';

export default {
  computed: {
    sortOptions() {
      return sortOptions;
    },

    sortOption: {
      get() {
        return this.sortOptions.find(option => option.key === this.value) || this.sortOptions[0];
      },

      set(value) {
        this.$emit('input', value.key);
      },
    },
  },

  props: {
    value: {
      type: String,
      required: false,
      default: () => sortOptions[0].key,
    },
    fluid: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
select {
  text-align: right;
  appearance: none;
  color: white;
  border: 0;
  line-height: 1;
  outline: 0;
  background-color: transparent;
  font-size: 1em;
  color: rgba(255,255,255,0.6);

  option {
    color: black;
    background-color: white;
  }

  &.fluid {
    width: 100%;
  }
}

// Chrome (at least the mobile version) does not support
// right-aligned selects or select options as of 2019-10-01.
// These worked fine in Firefox and I assumed it worked everywhere.

// This is a semi-hack to force right-aligned selects:
select {
  direction: rtl;
  option {
    direction: ltr;
  }
}
</style>
