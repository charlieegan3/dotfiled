<chunk>
  <div>
    <div class="file-chunk clearfix">
		<p class="contents">{ chunk.Contents }</p>
		<p class="details">
		  <span class="file-type { chunk.FileType }">{ chunk.FileType }</span>
		  <span class="count">used <strong>{ chunk.Files.length }</strong> times</span>
		</p>
	</div>
  </div>

  <li class="chunk-file clearfix" each={ chunk.Files }>
    <a href="{Repo}/find/master">{prettyRepoString(Repo)}</a> <span>{Name}</span>
  </li>

  <script>
    this.chunk = opts.chunk
    this.prettyRepoString = function(repoUrl) {
      parts = repoUrl.split("/");
      return parts.slice(parts.length - 2, parts.length).join("/");
    }
  </script>
</chunk>
