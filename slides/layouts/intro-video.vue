<script setup lang="ts">
import { ref } from "vue";

/**
 * Layout: intro-video
 * A layout with a full-screen background video and custom styled controls similar to AudioPlayer.
 * Controls are visible only on hover at the bottom center. Content is absolute.
 */
defineProps<{ video?: string }>();

const videoRef = ref<HTMLVideoElement>();
const isPlaying = ref(false);
const time = ref(0);
const duration = ref(0);

const fmt = (t: number) =>
    `${Math.floor(t / 60)}:${Math.floor(t % 60)
        .toString()
        .padStart(2, "0")}`;

function toggle() {
    if (!videoRef.value) return;
    if (isPlaying.value) videoRef.value.pause();
    else videoRef.value.play();
    isPlaying.value = !isPlaying.value;
}
</script>

<template>
    <div
        class="slidev-layout w-full h-full bg-black relative overflow-hidden group"
    >
        <!-- Background Video -->
        <video
            ref="videoRef"
            v-if="video"
            :src="video"
            class="absolute inset-0 w-full h-full object-cover z-0"
            @timeupdate="time = videoRef?.currentTime || 0"
            @loadedmetadata="duration = videoRef?.duration || 0"
            @ended="isPlaying = false"
        />

        <!-- Content Slot - Kept absolute for custom positioning in markdown -->
        <div class="absolute inset-0 z-10 pointer-events-none p-10">
            <slot />
        </div>

        <!-- Custom Controls - Absolute Bottom Center -->
        <div
            class="absolute bottom-6 left-1/2 -translate-x-1/2 z-20 flex items-center gap-2 bg-black/60 backdrop-blur-sm border border-white/10 p-1.5 rounded w-48 text-white select-none pointer-events-auto shadow-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"
        >
            <button
                @click="toggle"
                class="focus:outline-none hover:text-lemonYellow transition-colors flex-shrink-0"
            >
                <div
                    :class="
                        isPlaying
                            ? 'i-carbon:pause-filled'
                            : 'i-carbon:play-filled'
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
                        videoRef!.currentTime = Number(
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
    </div>
</template>

<style scoped>
/* Ensure content in slots remains interactive */
.slidev-layout :deep(div),
.slidev-layout :deep(h1),
.slidev-layout :deep(p),
.slidev-layout :deep(span),
.slidev-layout :deep(button),
.slidev-layout :deep(a) {
    pointer-events: auto;
}

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
