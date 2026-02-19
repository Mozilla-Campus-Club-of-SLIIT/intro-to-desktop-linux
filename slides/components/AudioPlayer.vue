<script setup lang="ts">
import { ref } from "vue";
defineProps<{ src: string }>();
const audio = ref<HTMLAudioElement>(),
    isPlaying = ref(false),
    time = ref(0),
    duration = ref(0);
const fmt = (t: number) =>
    `${Math.floor(t / 60)}:${Math.floor(t % 60)
        .toString()
        .padStart(2, "0")}`;
</script>

<template>
    <div
        class="flex items-center gap-2 bg-black border border-white/10 p-1.5 rounded w-48 mx-auto text-white select-none"
    >
        <audio
            ref="audio"
            :src="src"
            @timeupdate="time = audio?.currentTime || 0"
            @loadedmetadata="duration = audio?.duration || 0"
            @ended="isPlaying = false"
        />

        <button
            @click="
                isPlaying ? audio?.pause() : audio?.play();
                isPlaying = !isPlaying;
            "
            class="focus:outline-none hover:text-lemonYellow"
        >
            <div
                :class="
                    isPlaying ? 'i-carbon:pause-filled' : 'i-carbon:play-filled'
                "
                class="text-lg"
            />
        </button>

        <div class="flex-grow flex flex-col">
            <input
                type="range"
                min="0"
                :max="duration"
                step="0.1"
                :value="time"
                @input="
                    audio!.currentTime = Number(
                        ($event.target as HTMLInputElement).value,
                    )
                "
                class="w-full h-1 bg-white/20 accent-white cursor-pointer"
            />
            <div
                class="flex justify-between text-[8px] opacity-50 font-mono mt-0.5"
            >
                <span>{{ fmt(time) }}</span>
                <span>{{ fmt(duration) }}</span>
            </div>
        </div>
    </div>
</template>

<style scoped>
input[type="range"] {
    appearance: none;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 1px;
}
input[type="range"]::-webkit-slider-thumb {
    appearance: none;
    height: 6px;
    width: 6px;
    border-radius: 50%;
    background: white;
}
input[type="range"]::-moz-range-thumb {
    height: 6px;
    width: 6px;
    border-radius: 50%;
    background: white;
    border: none;
}
</style>
