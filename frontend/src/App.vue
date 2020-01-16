<template>
  <div id="app">

    <div id="createButtons">
      <!-- <vue-button type="default" v-on:click="showInvokeModal({'function': 'init_product_listing', 'fields': ['Product Listing Id', 'Supplier ID', 'Product ID'], 'title': 'Create Product Listing'})">Create Product Listing</vue-button>
      <vue-button type="default" v-on:click="showInvokeModal({'function': 'init_product', 'fields': ['Product Id', 'Quantity', 'CountryId'], 'title': 'Create Product'})">Create Product</vue-button> -->
      <vue-button type="default" v-on:click="showInvokeModal({'function': 'init_user', 'fields': ['ID', 'Company'], 'title': 'Create User'})">Create User</vue-button>
      <!-- <vue-button type="default" v-on:click="showInvokeModal({'function': 'init_regulator', 'fields': ['ID'], 'title': 'Create Regulator'})">Create Regulator</vue-button> -->
      <!-- <vue-button type="default" v-on:click="showInvokeModal({'function': 'update_work_order', 'fields': ['Work Order Number', 'New Owner ID'], 'title': 'Update Work Order'})">Update Work Order</vue-button> -->
      <!-- <vue-button type="default" v-on:click="showSelectModal">Update Work Order</vue-button> -->
      <vue-button type="default" v-on:click="getLedger">Refresh Ledger</vue-button>
    </div>
    </br>

    <!-- <div id="invokeQueryButtons">
      <vue-button type="default" v-on:click="showInvokeModal({'function': 'transfer_product_listing', 'fields': ['Product Listing Id', 'New Owner ID'], 'title': 'Transfer Product Listing'})">Transfer Product Listing</vue-button>
      <vue-button type="default" v-on:click="showInvokeModal({'function': 'check_products', 'fields': ['Product Listing ID', 'Regulator ID'], 'title': 'Check Products'})">Check Products</vue-button>
      <vue-button type="default" v-on:click="showInvokeModal({'function': 'update_work_order', 'fields': ['Word Order Number', 'New Owner ID'], 'title': 'Update Work Order'})">Update Work Order</vue-button>
      <vue-button type="default" v-on:click="getLedger">Refresh Ledger</vue-button>
    </div> -->

    <!-- <button
      class="btn green"
      @click="$modal.show('demo-login')">
      Demo: Login
    </button> -->
    <!-- <p>Ledger State -->
      <!-- <pre>
      {{ ledgerState | pretty }}
      </pre> -->
    <!-- </p> -->

    <!-- <data-table :data="gridData">
    </data-table> -->
    <!-- <v-card> -->


    <div>

      <modal name="update-wo-modal" height="auto" @before-open="beforeOpen">
        <vue-form>
          <div style="margin:auto;margin-left:120px">
              <!-- <slot name="body"> -->
                <h3 >{{modal_config.message}}</h3>
              <!-- </slot> -->
          </div>

          <template v-for="(field, idx) in modal_config.fields">
            <vue-form-item style="width:500px;" align=center>
              <vue-input
                :placeholder=field
                v-model=input[idx]>
              </vue-input>
            </vue-form-item>
          </template>


          <!-- <template v-if="modal_config.user_selections">
            <p>
              Select User
            </p>
            <select id="wo_user">
              <template v-for="v in modal_config.user_selections" >
              <option value=v>
                {{v}}
              </option>
            </template>
            </select>
          </template> -->

          <!-- <template v-if="modal_config.priority_selections">
            <p>
              Select Work Priority
            </p>
            <select id="wo_priority">
              <template v-for="p in modal_config.priority_selections" >
              <option value=p>
                {{p}}
              </option>
            </template>
            </select>
          </template> -->

          <vue-form-item style="margin-left:100px">
            <vue-button type="default" v-on:click="hideModal('update-wo-modal')">Cancel</vue-button>
            <vue-button type="success" v-on:click="invoke() ; hideModal('update-wo-modal')">Submit</vue-button>
          </vue-form-item>
        </vue-form>
        <!-- <template v-for="(field, idx) in fields">
          <vue-form-item style="width:500px;" align=center>
            <vue-input
              :placeholder=field
              v-model=input[idx]>
            </vue-input>
          </vue-form-item>
        </template> -->

        <!-- TODO add different modals based on status -->
        <!-- Show status -->
        <!-- template v-if "status == 'WAPPR'" -->
          <!-- Work order has been initiated via Maximo User, select a hazmat inspector to determine <scope of> work -->
          <!-- show hazmat inspectors -->
          <!-- allow dropdown to select priority/danger -->
        <!-- v-else-if == 'APPR' -->
          <!-- hazmat has verified work needs to be done -->
          <!-- show hazmat vendors -->
          <!-- vendor has already been assigned via Maximo -->
          <!-- click "Complete" if ready for work to be inspected -->
        <!-- v-else-if == 'PROGRESS' -->

      </modal>

      <modal name="invoke-modal" height="auto" >
        <h2 align="center"> {{title}} </h2>
        <vue-form
          id="chaincode-form"
          :model="form">

          <template v-for="(field, idx) in fields">
            <vue-form-item style="width:500px;" align=center>
              <vue-input
                :placeholder=field
                v-model="input[idx]">
              </vue-input>
            </vue-form-item>
          </template>

          <template v-if="func == 'init_user'">
            <vue-form-item >
              <v-select @change="appendType" style="width:400px;margin-left:100px" id="user_type" v-model="user_type" placeholder="User Type" :options="['HAZMAT_VENDOR','HAZMAT_INSPECTOR', 'BUILDING_INSPECTOR']"></v-select>
            </vue-form-item>
          </template>

          <vue-form-item style="margin-left:100px">
            <vue-button type="default" v-on:click="hideModal('invoke-modal')">Cancel</vue-button>
            <vue-button type="success" v-on:click="invoke">Submit</vue-button>
          </vue-form-item>
          </vue-form>
      </modal>

      <modal name="history-modal" height="auto" width="1000px" scrollable="true">
        <h2 align="center"> {{title}} </h2>
            <vuetable ref="vuetable"
              :api-mode="false"
              :data="history"
              :fields="historyfields"

            >
            </vuetable>
            <vue-button type="default" v-on:click="hideModal('history-modal')">Cancel</vue-button>
      </modal>

      <modal name="select-modal" height="auto" >
        <h2 align="center"> Select Work Order </h2>
          <!-- {{ledgerState.workorders}} -->
          <vue-form
            id="chaincode-form"
            :model="form">
            <vue-form-item style="margin-left:100px">

          <select v-model="selected">
            <option v-for="option in ledgerState.workorders" v-bind:value="option.id">
              {{ option.id }}
            </option>
          </select>
          <span>Selected: {{ selected }}</span>
        </vue-form-item>

          <vue-form-item style="margin-left:100px">
            <vue-button type="default" v-on:click="hideModal('select-modal')">Cancel</vue-button>
            <vue-button type="success" v-on:click="hideModal('select-modal') ; showInvokeModal({'function': 'update_work_order', 'fields': ['Work Order Number', 'New Owner ID'], 'title': 'Update Work Order'})">Next</vue-button>
          </vue-form-item>
          </vue-form>
      </modal>

      <div>

        <template v-if="ledgerState.workorders">
          <div>


          <h3 class="ui header">Work Orders</h3>
          <vuetable ref="vuetable"
            :api-mode="false"
            :data="ledgerState.workorders"
            :fields="['id', 'asset', 'vendor', 'status', 'lastmodifiedby', 'actions']"
          >
          <template slot="actions" scope="props">
            <div >
              <vue-button type="default"

                v-on:click="showModal('update-wo-modal', props.rowData)">
                <!-- <i class="zoom icon"></i> -->
                Update
              </vue-button>
              <vue-button type="default"
                v-on:click="getHistory(props.rowData)">
                View History
                <!-- TODO -->
                <!-- https://stackoverflow.com/questions/51136060/how-to-get-transaction-history-on-a-particular-key-in-hyperledger-fabric -->
              </vue-button>
              <!-- TODO, add button to view ledger history for WO -->
            </div>
          </template>
        </vuetable>
        </div>

            <!-- <data-table :data=ledgerState.workorders :style="{width: '800px', height: '800px', overflow: 'auto'}">
              <template slot="caption">Work Orders</template>
            </data-table> -->
        </template>

        <template v-if="ledgerState.users">
                <!-- <data-table :data="ledgerState.users" :style="{width: '300px', height: '200px', overflow: 'auto'}">
                    <template slot="caption">Users</template>
                </data-table> -->
                <vuetable ref="vuetable"
                  :api-mode="false"
                  :data="ledgerState.users"
                  :fields="Object.keys(ledgerState.users[0])"
                >
                </vuetable>
        </template>

        <template v-if="ledgerState.suppliers">
                <data-table :data="ledgerState.suppliers.map( s =>  ({ Id: s.Id, CountryId: s.countryId, OrgId: s.orgId }) )" :style="{width: '400px', height: '200px', overflow: 'auto'}">>
                  <template slot="caption">Suppliers</template>
                </data-table>
        </template>
      </div>
      <!-- </template> -->
      <!-- </template> -->
      </div>
      <!-- </v-card> -->
      <!-- <p>Ledger State</p> -->
      <div class="col-xs-12">
        <div class="well">
          <tree-view :data="ledgerState" :options="{maxDepth: 1, rootObjectKey: 'ledgerState'}"></tree-view>
        </div>
      </div>
      </br>
      </br>
      </br>
      </br>

<!--  TODO, this is being hidden by datatables since it's fixed -->
<!-- <div v-if="!isHidden" style="z-index:9000">
      <vue-form
        id="chaincode-form"
        :model="form"
        style="width: 75%; position: fixed; left: 50%; margin-left: -37.5%;">
        <h2 style="float:center"> Invoke Chaincode </h2>
        <vue-form-item label="Function">
          <vue-input
            placeholder="Function"
            v-model="form.function"
            style="width: 100%">
          </vue-input>
        </vue-form-item>

        <vue-form-item label="Arguments">
          <vue-input
            placeholder="Arguments"
            v-model="form.args"
            style="width: 100%">
          </vue-input>
        </vue-form-item>
        <vue-form-item>
          <vue-button type="default" v-on:click="isHidden = true">Cancel</vue-button>
          <vue-button type="success" v-on:click="invoke" >Submit</vue-button>
        </vue-form-item>
      </vue-form>
    </div> -->
<!-- <vue-form-item> item 1 </vue-form-item>
      <vue-form-item> item 2 </vue-form-item> -->
<!-- <vue-input placeholder="Please input"></vue-input>
      <vue-input placeholder="Please input"></vue-input> -->



</div>


</template>

<script>
  import 'vfc/dist/vfc.css'
  import './dist/json-tree.css'

  import {
    Input
  } from 'vfc'
  import {
    Form
  } from 'vfc'
  import {
    FormItem
  } from 'vfc'
  import {
    Button
  } from 'vfc'
  // import DemoLoginModal       from './components/modals/DemoLoginModal.vue'

  // import 'vue-js-modal'
  // import { Card } from 'v-card'
  // import { DataTable } from 'v-data-table'
  // import { Button } from 'vfc'



  export default {
    name: 'app',
    data() {
      return {
        isHidden: false,
        form: {
          function: '',
          args: ''
        },
        args: [],
        ledgerState: {},
        products: [],
        fields: [],
        user_fields: [],
        user_type: '',
        status: '',
        user_input: [],
        modal_config: {},
        selected_wo: '',
        input: [],
        func: '',
        history: [],
        historyfields: [],
        title: ''
        // gridData: [{"id":"product1","quantity":"300","countryId":"US"},{"id":"product2","quantity":"350","countryId":"US"}]
      }
    },


    mounted () {
      this.getLedger()
    },
    components: {
      Form,
      FormItem,
      // DemoLoginModal,
      [Input.name]: Input,
      [Button.name]: Button
    },
    methods: {
      invoke() {
        console.log("this.$data.input")
        console.log(this.$data.input)
        // console.log(this.$data.user_input)
        // console.log(this.$data.input.concat(this.$data.user_input))
        if (this.$data.func == 'update_work_order') {
          // this.$data.input.splice(1, 0, this.$data.status).splice(1, 0, this.$data.selected_wo)
          var args = [this.$data.selected_wo, this.$data.status].concat( this.$data.input )
        } else {
          var args = this.$data.input
        }
        console.log("args")
        console.log(args)
        var options = {
          method: "POST",
          body: JSON.stringify({
            method: "invoke",
            params: {
              ctorMsg: {
                function: this.$data.func, //this.$data.input[0],
                args: args //this.$data.input
                // function: this.$data.form.function,
                // args: this.$data.form.args.split(',') // ["retailer1", "retailer"]
              }
            }
          }),
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          }
        }
        console.log(options)
        fetch("http://localhost:3000/api/chaincode", options).then(() => {
          console.log("api call complete")
          // this.$data.isHidden = true
          this.$modal.hide('invoke-modal');
          this.$data.input = []
          this.$data.user_input = []
          // auto refresh ledger after 3 seconds
          setTimeout( () => this.getLedger(), 3000 )
        })
      },
      getHistoryFields() {
        var fields = ['txId'].concat(Object.keys(this.$data.history[0].value))
        this.$data.historyfields = fields
        // return fields
      },
      getHistory(rowData) {
        this.$data.history = []
        var options = {
          method: "POST",
          body: JSON.stringify({
            method: "query",
            params: {
              ctorMsg: {
                function: "getHistory",
                args: [rowData.id]
              }
            }
          }),
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          }
        }
        fetch("http://localhost:3000/api/chaincode", options).then((response) => {
          console.log("ledger history retrieved")
          response.json().then((json) => {
            console.log(json)
            var result = JSON.parse(json)
            result.map( (r) => {
              this.$data.history.push( Object.assign({txId: r.txId.substring(1, 10) + '...'}, r.value ) )
            })

            setTimeout( () => this.$modal.show('history-modal'), 1000 )
            this.$data.historyfields = ['txId'].concat(Object.keys(result[0].value))
            // this.showHistoryModal()
            // history-modal

            // this.$data.ledgerHistory

          })

        })
      },
      getLedger() {
        console.log(this.$data.form.function)
        console.log("publishing")

        var options = {
          method: "POST",
          body: JSON.stringify({
            method: "query",
            params: {
              ctorMsg: {
                function: "read_everything",
                args: []
              }
            }
          }),
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          }
        }
        fetch("http://localhost:3000/api/chaincode", options).then((response) => {
          console.log("ledger state retrieved")
          response.json().then((json) => {
            console.log(json)
            var result = JSON.parse(json)
            // for (reg in json.regulators) {
            //   delete reg.id
            // }
            // for (reg in json.retailers) {
            //   delete reg.id
            // }

            this.$data.ledgerState = result
            // this.$data.products = JSON.parse(json.products)
            // console.log("json.products")
            // console.log(json.products)
          })

        })
      },

      /*
      <!-- TODO add different modals based on status -->
      <!-- Show status -->
      <!-- template v-if "status == 'WAPPR'" -->
        <!-- Work order has been initiated via Maximo User, select a hazmat inspector to determine <scope of> work -->
        <!-- show hazmat inspectors -->
        <!-- allow dropdown to select priority/danger -->
      <!-- v-else-if == 'APPR' -->
        <!-- hazmat has verified work needs to be done -->
        <!-- show hazmat vendors -->
        <!-- vendor has already been assigned via Maximo -->
        <!-- click "Complete" if ready for work to be inspected -->
      <!-- v-else-if == 'PROGRESS' -->
      */
      // INPRG, COMP, WAPPR, APPR, COMP,
      // showHistoryModal() {
      //   var name = 'history-modal'
      //   this.$modal.show(name);
      // },
      showModal(name, rowData) {
        console.log(`opening modal ${name}`)
        this.$data.input = []
        if (name == 'update-wo-modal') {
          this.$data.func = "update_work_order"
          this.$data.selected_wo = rowData.id
          if (rowData.status == 'WAPPR') { // 'WAPPR') {
            this.$data.status = rowData.status
            var userType = "HAZMAT_INSPECTOR"
            // var message = `Work order issued by Maximo User. ${userType}, please enter user ID and work priority (1-4) below` // TODO, may want to add actual user name at some point
            var message = `Enter ${userType} user ID and work priority (1-4) below`
            var users = this.$data.ledgerState.users.filter( (u) => u.type == userType ).map( (m) => m.id )
            var priority = [1,2,3,4] // 4 is emergency
            var id = rowData.id
            var fields = ["User ID", "Work Priority"]
            var config = {
              "id": id,
              "title": `Work Order ${id}`,
              "text": message,
              "message": message,
              "status": rowData.status,
              "user_selections": users, // filtered by type
              "priority_selections": priority,
              "user_type": userType,
              "fields": fields
            }
            this.$data.modal_config = config
          } else if (rowData.status == 'APPR') {
            this.$data.status = rowData.status
            // var message = `Hazmat Inspector ${rowData.lastmodifiedby} has verified work that needs to be done. Vendor: ${rowData.vendor} has already been assigned via Maximo. If vendor work is complete and ready for inspection, please enter Vendor's user ID and press Submit. `
            var message = `Enter Hazmat Vendor's user ID and press Submit. `
            var users = null
            var id = rowData.id
            var fields = ["User ID"]
            // this.$data.input[0] = rowData.vendor
            var config = {
              "id": id,
              "message":  message,
              "status": rowData.status,
              // "user_selections": users
              "fields": fields
            }
            this.$data.modal_config = config
          } else if (rowData.status == 'INPRG') {
            this.$data.status = rowData.status
            // var message = `Vendor ${rowData.vendor} has already been assigned via Maximo. If ready for work to be inspected, enter building inspector ID and click \"Complete\" `
            var message = `Enter building inspector ID and click \"Complete\" `
            var userType = 'BUILDING_INSPECTOR'
            var users = this.$data.ledgerState.users.filter((u) => u.type == userType )
            var id = rowData.id
            var fields = ["User ID"]
            var config = {
              "id": id,
              "message":  message,
              "status": rowData.status,
              "user_selections": users,
              "user_type": userType,
              "fields": fields
            }
            this.$data.modal_config = config
            // var fields =
          }
        }
        this.$modal.show(name, config);
      },
      beforeOpen (event) {
        console.log(event.params);
        this.modal_config = event.params
      },
      showInvokeModal(config) {
        console.log("opening modal")
        // console.log(fields)
        this.$data.input = []
        console.log(config.function)
        console.log(config.fields)
        this.$data.func = config.function
        this.$data.fields = config.fields
        this.$data.title = config.title
        this.$data.user_fields = []
        this.$data.user_type = ''
        this.$modal.show('invoke-modal', {
          "fields": config.fields
        });
      },
      hideModal(name) {
        this.$modal.hide(name);
        console.log("hiding modal")
        this.$data.user_fields = []
        this.$data.user_type = ''
      },
      getFormValues() {
        console.log("getting form vals")
        console.log(this.$data.input)
        // this.output = this.$refs.my_input.value
        // console.log(this.$refs.my_input.value)
      },
      appendType() {
      //   this.$data.input.push(this.$data.user_type)
        console.log("this.user_type")
        console.log(this.user_type)
        console.log("this.$data.input")
        console.log(this.$data.input)
        this.$data.input[2] = this.user_type
        // if (this.user_type == "supplier") {
      //     this.$data.user_fields.push('Country Id')
      //     this.$data.user_fields.push('Org Id')
      //   } else {
      //     this.$data.user_fields = []
      //   }
      }

    },
    filters: {
      pretty: function(value) {
        return JSON.stringify(value, null, 2);
      }
    }


  }
</script>

<!-- TODO, finish modal -->
<!-- <script type="text/x-template" id="modal-template">
  <transition name="modal">
    <div class="modal-mask">
      <div class="modal-wrapper">
        <div class="modal-container">

          <div class="modal-header">
            <slot name="header">
              default header
            </slot>
          </div>

          <div class="modal-body">
            <slot name="body">
              default body
            </slot>
          </div>

          <div class="modal-footer">
            <slot name="footer">
              default footer
              <button class="modal-default-button" @click="$emit('close')">
                OK
              </button>
            </slot>
          </div>
        </div>
      </div>
    </div>
  </transition>
</script> -->


<style>
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    margin-top: 60px;
  }
</style>
