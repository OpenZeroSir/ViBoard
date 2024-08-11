import { createStore } from 'vuex'
export interface PublicState {
    //全局websocket
    websocket: WebSocket|null,
    //后端生成的 websocket token
    wstoken:string
 }
export const publicStore = createStore<PublicState>({
    state: {
        websocket: null,
        wstoken:"",
    },
    mutations: {
        set_wssock(state, s: WebSocket) {
            state.websocket = s
        },
        set_wstoken(state,t:string){
            state.wstoken = t
        },
    }
})