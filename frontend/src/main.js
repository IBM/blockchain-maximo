import Vue from 'vue'
import App from './App.vue'
import VFC from 'vfc'
import DataTable from 'v-data-table'
import VModal from 'vue-js-modal'
// import Modal from './components/Modal.vue'

import vSelect from 'vue-select'

import 'vfc/dist/vfc.css'
import './dist/json-tree.css'

import Vuetable from 'vuetable-2/src/components/Vuetable'
import VuetablePagination from 'vuetable-2/src/components/VuetablePagination'
import VuetablePaginationInfo from 'vuetable-2/src/components/VuetablePaginationInfo'

import TreeView from "vue-json-tree-view"

Vue.config.productionTip = false

Vue.use(DataTable)
Vue.use(VFC)
Vue.use(VModal)
Vue.use(TreeView)
Vue.component('v-select', vSelect)
// Vue.use(VModal, { componentName: "modal" })
Vue.component("vuetable", Vuetable);
Vue.component("vuetable-pagination", VuetablePagination);
// Vue.component("vuetable-pagination-dropdown", VuetablePaginationDropDown);
Vue.component("vuetable-pagination-info", VuetablePaginationInfo);

// Vue.use(Vuetable)
// Vue.component('modal', {
//   template: '#modal-template'
// })
import CssForBootstrap4 from './components/VuetableCssBootstrap4.js'
new Vue({
  render: h => h(App),
  data: function() {
    return {
      showModal: false
      // css: CssForBootstrap4
    }
  }
}).$mount('#app')
