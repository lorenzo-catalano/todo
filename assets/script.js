Vue.component('task-list', {
  template: `<div>
      <input v-model="newTask" type="text"/> <button @click="add">add</button>
      <ul id="example-1">
      <li v-for="item in tasks" :key="item.ID">
          <input type="checkbox" v-model="item.completed" @change="complete(item.ID,$event.target.checked)"/>
          <span class="grow1" @click="select(item.ID)">{{ item.task }}</span>
          <span class="action" title="Delete" @click="remove(item.ID)"></span>
      </li>
      </ul>
  </div>`,
  data: function () {
    return {
      tasks: [],
      newTask:''
    }
  },
  created : function(){
      this.loadData()
  },
  methods: {
    loadData: async function () {
          var tasks = await fetch('/tasks').then(d=>d.json())
          this.tasks = tasks
    },
    complete: async function(taskId,completed){
      await fetch(`/tasks/${taskId}`,{
          method:'POST',
          body: JSON.stringify({completed:completed})
      }).then(d=>d.json())
      this.loadData()
    }, 
    add: async function(){
      var nrn = this.newTask
      this.newTask = ''
      await fetch('/tasks',{
          method:'PUT',
          body: JSON.stringify({task:nrn})
      }).then(d=>d.json())
      this.loadData()
    }, 
    remove: async function(taskId){
      await fetch(`/tasks/${taskId}`,{method:'DELETE'}).then(d=>d.json())
      this.loadData()
    }
  }
})
  
var app = new Vue({
    el: '#app'
})