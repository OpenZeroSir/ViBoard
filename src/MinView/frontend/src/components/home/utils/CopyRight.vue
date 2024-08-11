<template>
    <v-dialog v-model="store.state.viShowCopyRight" width="auto">
        <v-card>
            <div style="background-color: white;padding:10px;">
                <img src="../../../assets/images/home/home.webp" style="width: 100%;"> 
            </div>
            <div style="background-color: white;padding:10px;">
                <img src="../../../assets/images/home/home1.webp" style="width: 100%;">
            </div>
            <div style="background-color: white;padding:10px;">
                <img src="../../../assets/images/home/home2.webp" style="width: 100%;">
            </div>
            <v-card-text>
                <div style="text-align: center;">
                    Copyright © 2024 zerosir.com
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { store } from '../../../store/index';


//这里是编辑器的版权

</script>