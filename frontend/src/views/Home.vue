<template>
  <div class="container-fluid" style="padding-left: 15px !important; padding-right: 15px !important;">
    <div class="row justify-content-center">
      <div class="col-12 col-lg-12 col-xl-12 content-pop mt-0">

        <div class="container-fluid">

          <div class="row" id="changelog">
            <div class="col-12">
              <v-runtime-template class="changelog markdown-body" :template="changelog"/>
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script>

import EqWindow        from "@/components/eq-ui/EQWindow";
import UserContext     from "@/app/user/UserContext";
import {SpireApi}      from "../app/api/spire-api";
import * as util       from "util";
import VideoViewer     from "../app/video-viewer/video-viewer";
import LazyImageLoader from "@/app/lazy-image-load/lazy-image-load";

export default {
  components: {
    EqWindow,
    "v-runtime-template": () => import("v-runtime-template")
  },
  data() {
    return {
      userContext: null,
      changelog: "",
    }
  },
  async mounted() {
    this.userContext = await (UserContext.getUser())

    SpireApi.v1().get(`/app/changelog`).then((response) => {
      if (response.data && response.data.data) {

        let markdownRaw = response.data.data

        const youTubeSplit = markdownRaw.split("[![](https://img.youtube.com/vi/")

        youTubeSplit.forEach((e) => {
          if (e.includes("/0.jpg)](https://www.youtube.com")) {
            const videoCodeSplit = e.split("/0.jpg")
            if (videoCodeSplit.length > 0) {
              const videoCode = videoCodeSplit[0].trim()

              // replace markdown code for html
              markdownRaw = markdownRaw.replace(
                util.format("[![](https://img.youtube.com/vi/%s/0.jpg)](https://www.youtube.com/watch?v=%s)", videoCode, videoCode),
                util.format('<div class="container"><iframe allow="autoplay" class="video" src="https://www.youtube.com/embed/%s?mute=1&showinfo=0&controls=0&modestbranding=1&rel=0&loop=1&showsearch=0&iv_load_policy=3&playlist=%s" title="YouTube video player" frameborder="0" allowfullscreen></iframe></div>\n', videoCode, videoCode)
              )

              // console.log("Video code is [%s]", videoCode)
            }
          }
        })

        const md = require("markdown-it")({
          html: true,
          xhtmlOut: false,
          breaks: true,
          typographer: false,
          linkify: true
        });

        markdownRaw = md.render(markdownRaw);

        // lazy image load injection
        markdownRaw = markdownRaw.replaceAll(
          "img src=",
          "img class='lazy-image lazy-image-unloaded' src='data:image/gif;base64,R0lGODlhAQABAIAAAMLCwgAAACH5BAAAAAAALAAAAAABAAEAAAICRAEAOw==' data-src="
        )

        // doc
        this.changelog = "<div>" + markdownRaw + "</div>"

        setTimeout(() => {
          const anchors = document.getElementById('changelog').getElementsByTagName('a');
          for (var i = 0; i < anchors.length; i++) {
            anchors[i].setAttribute('target', '_blank');
          }

          document.querySelectorAll('#changelog h1, #changelog h2, #changelog h3, #changelog h4').forEach($heading => {

            //create id from heading text
            const id = $heading.getAttribute("id") || $heading.innerText.toLowerCase().replace(/[`~!@#$%^&*()_|+\-=?;:'",.<>\{\}\[\]\\\/]/gi, '').replace(/ +/g, '-');

            //add id to heading
            $heading.setAttribute('id', id);

            //append parent class to heading
            $heading.classList.add('anchor-heading');

            //create anchor
            let $anchor         = document.createElement('a');
            $anchor.className   = 'anchor-link';
            $anchor.href        = '#' + id;
            $anchor.innerText   = ' # ';
            $anchor.style.color = '#666';

            //append anchor after heading text
            $heading.append($anchor);
          });

          document.querySelectorAll("table").forEach((e) => {
            if (e) {
              e.classList.add('eq-table')
              e.classList.add('bordered')
              e.outerHTML = "<div class='eq-window-simple mt-3 mb-3 p-0' style='overflow-y: hidden'>" + e.outerHTML + "</div>"
            }
          });

        }, 100)

      }
    })

    LazyImageLoader.addScrollListener()

    // auto play videos that are in the viewport
    window.addEventListener("scroll", this.handleRender);
    setTimeout(() => {
      this.handleRender()
      LazyImageLoader.handleRender()
    }, 500)
  },
  methods: {
    handleRender() {
      let videos = document.getElementsByClassName("video");
      for (let i = 0; i < videos.length; i++) {
        let video = videos.item(i)
        if (VideoViewer.elementInViewport(video) && !video.src.includes("autoplay")) {
          video.src = video.src + "&autoplay=1"
        }
      }
    }
  },
  deactivated() {
    window.removeEventListener("scroll", this.handleRender, false)
    LazyImageLoader.destroyScrollListener()
  }
}
</script>

<style>
.changelog {
  font-size: 16px;
  line-height: 1.5;
  word-wrap: break-word;
  -webkit-text-size-adjust: 100%;
}

.changelog img {
  border-radius: 5px;
  text-align: center;
  display: block;
  margin-bottom: 5px;
}

.container {
  position: relative;
  width: 100%;
  height: 0;
  padding-bottom: 56.25%;
}

.video {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  padding-bottom: 30px;
  padding-top: 15px;
}

.changelog code {
  padding: 0.2em 0.4em;
  margin: 0;
  color: red;
  background-color: rgba(110, 118, 129, .4);
  border-radius: 6px;
}

.markdown-body h1 {
  padding-bottom: 0.6em;
  font-size: 2em;
  border-bottom: 1px solid #6666665e;
}

.markdown-body h2 {
  font-size: 1.5em;
  padding-bottom: 0.6em;
  border-bottom: 1px solid #6666665e;
}

.markdown-body h3 {
  font-size: 1.25em;
}

.markdown-body h4 {
  font-size: 1em;
}

.markdown-body h1, .markdown-body h2, .markdown-body h3, .markdown-body h4, .markdown-body h5, .markdown-body h6 {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

.markdown-body img {
  max-width: 100%;
  box-sizing: content-box;
  background-color: black;
}

</style>
