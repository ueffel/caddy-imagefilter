# caddy-imagefilter

This package can transform images from the file system in various ways.

Supported formats are:

* JPEG
* GIF
* PNG
* BMP
* TIFF
* WEBP (only input)

There is currently no pure go WEBP library, that can encode (not just decode) WEBP images (maybe
[someday](https://github.com/golang/go/issues/45121)). Therefore images in WEBP will be encoded as
PNG as fallback.

It's recommended to keep transformed images in a cache to improve response times and don't do
transformation over and over again.

## Installation

### With Default filters

```sh
xcaddy build --with github.com/ueffel/caddy-imagefilter/defaults
```

See [Default filters](#default-filters).

### With All filters

```sh
xcaddy build --with github.com/ueffel/caddy-imagefilter/all
```

This includes [Default filters](#default-filters) as well as [Additional filters](#additional-filters)

### With Individual filters

```sh
xcaddy build --with github.com/ueffel/caddy-imagefilter/fit
```

or

```sh
xcaddy build --with github.com/ueffel/caddy-imagefilter/fit --with github.com/ueffel/caddy-imagefilter/resize
```

or any other combination.

## Configuration

### Ordering the directive

```caddy-d
{
    order image_filter before file_server
}
```

### Syntax

```caddy-d
image_filter [<matcher>] {
    root            <path>
    jpeg_quality    <quality>
    png_compression <level>
    max_concurrent  <level>

    # included filters
    <filters...> <filter-args...>
}
```

* **root** sets the path to the site root for just this file server instance, overriding any other.
  Default: `{http.vars.root}` or the current working directory. Note: This subdirective only changes
  the root for this directive. For other directives (like `try_files` or `templates`) to know the
  same site root, use the root directive, not this subdirective.
* **jpeg_quality** determines the quality of jpeg encoding after the filters are applied. It ranges
  from 1 to 100 inclusive, higher is better. Default is `75`.
* **png_compression** determines the compression of png images. Possible values are:
  * `0`: Default compression
  * `-1`: no compression
  * `-2`: fastest compression
  * `-3`: best compression
* **max_concurrent** determines how many request can be served concurrently. This is intended to
  reduce excessive cpu/memory usage for image transformations by limiting the number of parallel
  calculations. Any value less or equal `0` means no limit. Default is `0`.
* **<filters...>** is a list of filters with their corresponding arguments, that are applied in
  order of definition.
* **<filter-args...>** support [caddy
  placeholders](https://caddyserver.com/docs/caddyfile/concepts#placeholders).

At least one filter has to be configured or this would be just an inefficient `file_server`. No
filters configured is therefore considered invalid and will emit an error on start.

### Examples

```caddy-d
@thumbnail {
    path_regexp thumb /.+\.(jpg|jpeg|png|gif|bmp|tif|tiff|webp)$
    query w=*
    query h=*
}

image_filter @thumbnail {
    fit {query.w} {query.h}
}
```

Produces are scaled version of an image, that fits within a rectangle with width `w` and height `h`.
`w` and `h` are query parameters. For example: <http://localhost/image.png?w=400&h=225>

```caddy-d
@thumbnail {
    path_regexp thumb /w(400|800)(/.+\.(jpg|jpeg|png|gif|bmp|tif|tiff|webp))$
}
handle @thumbnail {
    rewrite {re.thumb.2}
    image_filter {
        resize {re.thumb.1} 0
    }
}
```

This configuration takes a parameter from the path and only allows a limited number of dimensions
(400 or 800) for resizing. For example: <http://localhost/w400/image.png>. It then resizes the image
down to the width while keeping the aspect ratio preserved (height is 0, see [resize](#resize)). The
rest of the path after `w400` is the file path to the image. (`/image.png`)

(The internal ordering of directives is important here (global `order` directive), `image_filter`
has to run after `rewrite`)

```caddy-d
image_filter {
    crop 500 500 topleft
    rotate 180
    flip v
    resize 200 400
    sharpen
}
```

You can go crazy and combine many filters. (But no more than 9999, which should quite sufficient or
you're doing something seriously wrong)

### Default filters

#### crop

Crop produces a cropped image as rectangular region of a specific size.

Syntax:

```caddy-d
    crop <width> <height> [<anchor>]
```

Parameters:

* **width** must be a positive integer and determines the width of the cropped image.
* **height** must be a positive integer and determines the height of the cropped image.
* **anchor** determines the anchor point of the rectangular region that is cut out. Possible values
  are: center, topleft, top, topright, left, right, bottomleft, bottom, bottomright. Default is
  center.

Installation: `--with github.com/ueffel/caddy-imagefilter/crop`

#### fit

Fit scales an image to fit to the specified maximum width and height using a linear filter, the
image aspect ratio is preserved. If the image already fits inside the bounds, nothing will be
done.

Syntax:

```caddy-d
    fit <width> <height>
```

Parameters:

* **width** must be a positive integer and determines the maximum width.
* **height** must be a positive integer and determines the maximum height.

Installation: `--with github.com/ueffel/caddy-imagefilter/fit`

#### flip

Flip flips (mirrors) a image vertically or horizontally.

Syntax:

```caddy-d
    flip <h|v>
```

Parameters:

* **h|v** determines if the image flipped horizontally or vertically.

Installation: `--with github.com/ueffel/caddy-imagefilter/flip`

#### resize

Resize can downsize images. If upsizing of an image is detected, nothing will be done and
the input image is returned unchanged.

Syntax:

```caddy-d
    resize <width> <height>
```

Parameters:

* **width** must be a positive integer and determines the maximum width.
* **height** must be a positive integer and determines the maximum height.

Either width or height can be 0, then the image aspect ratio is preserved.

Installation: `--with github.com/ueffel/caddy-imagefilter/resize`

#### rotate

Rotate rotates a image 90, 180 or 270 degrees counter-clockwise.

Syntax:

```caddy-d
    rotate <angle>
```

Parameters:

* **angle** is one of the following: 0, 90, 180, 270 (0 is valid, but nothing will be done to the
  image).

Installation: `--with github.com/ueffel/caddy-imagefilter/rotate`

#### sharpen

Sharpen produces a sharpened version of the image.

Syntax:

```caddy-d
    sharpen [<sigma>]
```

Parameters:

* **sigma** must be a positive floating point number and indicates how much the image will be
  sharpened. Default is 1.

Installation: `--with github.com/ueffel/caddy-imagefilter/sharpen`

### Additional filters

#### blur

Blur produces a blurred version of the image.

Syntax:

```caddy-d
    blur [<sigma>]
```

Parameters:

* **sigma** must be a positive floating point number and indicates how much the image will be
  blurred. Default is 1.

Installation: `--with github.com/ueffel/caddy-imagefilter/blur`

#### grayscale

Grayscale produces a grayscaled version of the image.

Syntax:

```caddy-d
    grayscale
```

no parameters.

Installation: `--with github.com/ueffel/caddy-imagefilter/grayscale`

#### invert

Invert produces an inverted (negated) version of the image.

Syntax:

```caddy-d
    invert
```

no parameters.

Installation: `--with github.com/ueffel/caddy-imagefilter/invert`

#### rotate_any

RotateAny rotates an image by a specific angle counter-clockwise. Uncovered areas after the rotation
are filled with the specified color.

Syntax:

```caddy-d
    rotate_any <angle> <color>
```

Parameters:

* **angle** is the angle as floating point number in degrees by which the image is rotated
  counter-clockwise.
* **color** is the color which is used to fill uncovered areas after the rotation. Supported formats
  are:
  * `"#FFAADD"` (in quotes because otherwise it will be a comment in a caddyfile)
  * `rgb(255,170,221)`
  * `rgba(255,170,221,0.5)`
  * `transparent`, `black`, `white`, `blue` or about 140 more

Installation: `--with github.com/ueffel/caddy-imagefilter/rotate_any`

#### smartcrop

Smartcrop finds good rectangular image crops of a specific size. It uses
<https://github.com/muesli/smartcrop>

Syntax:

```caddy-d
    smartcrop <width> <height>
```

Parameters:

* **width** must be a positive integer and determines the width of the cropped image.
* **height** must be a positive integer and determines the height of the cropped image.

Installation: `--with github.com/ueffel/caddy-imagefilter/smartcrop`

### Advanced Configuration

The following configuration uses <https://github.com/caddyserver/cache-handler> as internal cache,
so image transformations are done once per day and then served from the cache. Beware that this
cache module is not quite ready for production use, it should just serve as an example here.

```caddy-d
{
    order image_filter before file_server
    order cache first
    cache {
        olric_config olricd.yaml
    }
}

http://:80 {
    root .
    file_server browse
    @thumbnail {
        path_regexp thumb /w(400|800)/.+\.(jpg|jpeg|png|gif|bmp|tif|tiff|webp)$
    }
    reverse_proxy @thumbnail [::1]:9000 {
        header_up Cache-Control "max-age=86400"  # always use cache, even if requested otherwise
        header_down -Server                      # prevents multiple "Server: Caddy" response header
    }
}

http://:9000 { # internal address only accessable from the server itself to transform images
    bind ::1   # local ipv6 address (same as 127.0.0.1 for ipv4)
    cache
    header Cache-Control "max-age=86400" # keep 1 day in cache
    root .
    @thumbnail {
        path_regexp thumb (?i)/w([0-9]+)(/.+)$
    }
    handle @thumbnail {
        rewrite {re.thumb.2}
        image_filter {
            resize {re.thumb.1} 0
        }
    }
}
```

This could also extended to limit the load because of image filtering by using rate limiting with
this module <https://github.com/mholt/caddy-ratelimit>

## Write your own filter

You can use the base module `imagefilter` to implement your own filter. A new filter modules needs 2
types, that implement `imagefilter.FilterFactory` and `imagefilter.Filter` respectively. A
`FilterFactory` registers itself as image filter and produces configured `Filter` instances. The
configured `Filter` instance then is called at runtime with an image, where the filter operation has
to be applied and the resulting image returned. It's recommended to have all configure-parameters as
strings, so they can contain caddy placeholders. Before applying the filter the placeholders should
be replaced with `caddy.Replacer`'s `ReplaceAll`.

Have a look at the default filters for implementation pointers.
