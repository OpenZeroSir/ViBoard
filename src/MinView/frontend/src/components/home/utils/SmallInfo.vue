<template>
    <div style="display: flex;">
        <div style="display: flex;margin-top: 5px;">
            <img :src="wechat" alt="" :style="{
                height: '256px',
                marginRight: '10px',
                marginLeft:'10px',
            }">
            <img :src="alipay" alt="" style="height: 256px;">
        </div>
        <div style="margin-left: 20px;height: 265px;position: relative;">
            <div>
                <div :style="{
                    color: 'white',
                }">
                    我叫<span style="color: #ff9900;font-weight: bold;">李飞龙</span>，软件由我个人技术、爱好和资金在驱动，软件免费但无法保证任何包含但不仅限经济、安全和质量等问题，除了以下情况:
                </div>
                <div style="color: #ff9900;font-weight: bold;margin-top: 10px;">
                    1、软件不会植入广告
                </div>
                <div style="color: #ff9900;font-weight: bold;margin-bottom: 10px;">
                    2、软件不会留有后门
                </div>
                <div :style="{
                    color: 'white',
                }">
                    以上条件成立的前提是通过合规渠道获取到软件，任何形式交换获取软件或包含但不仅限参与软件官方或非官方活动等，均可以理解为骗局。
                </div>
            </div>
            <div style="bottom: 0; position: absolute;">
                <div :style="{
                    color: '#ff9900',
                    fontWeight:'bold',
                }">
                    如您发现软件问题，可以把详细问题发送到我邮箱，当然也可以请我喝杯咖啡。
                </div>
                <div style="color: white;">
                    邮箱：open@zerosir.com
                </div>
            </div>
        </div>
        
    </div>
</template>

<script setup lang="ts">
import alipay from "./../../../assets/images/pay/alipay.webp";
import wechat from "./../../../assets/images/pay/wechat.webp";

</script>