title = this blog
date = 2019-05-24
desc = building this blog

---

Handrolled, artisanal blog.
Lovingly made from scratch.
Decorated with hair torn out from trying to get things to work.

## Pipeline

### Source

Hosted on [Github](https://github.com/seankhliao/com-seankhliao-blog).
Written in markdown (mostly).

### CI/CD

Pushed commits trigger CI/CD on [Google Cloud Build](https://cloud.google.com/cloud-build/).
Runs custom site generator ([parse markdown](https://github.com/russross/blackfriday), create html).
Pushes results into [Firebase Hosting](https://firebase.google.com/products/hosting/).
Extra steps to purge Cloudflare cache.
<span>TODO: trigger [Wayback Machine](web.archive.org) archive and [Google Search Console](https://search.google.com/search-console/about) indexing</span>

### Front End

Dark theme, monospace font, actually sort of what my terminal looks like.
Mostly reused the stylesheet from my [main](https://seankhliao.com) site.

### Aditional Info

[kaniko](https://github.com/GoogleContainerTools/kaniko) is used to build containers.
Getting it to respect subdirectory contexts is still confusing
