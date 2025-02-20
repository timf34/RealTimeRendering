#iChannel0 "file://checkerboard.png"      // Checkerboard texture
#iChannel1 "file://Skyboxes/back.jpg"     // Skybox back
#iChannel2 "file://Skyboxes/front.jpg"    // Skybox front
#iChannel3 "file://Skyboxes/left.jpg"     // Skybox left
#iChannel4 "file://Skyboxes/right.jpg"    // Skybox right
#iChannel5 "file://Skyboxes/top.jpg"      // Skybox top
#iChannel6 "file://Skyboxes/bottom.jpg"   // Skybox bottom

const float PI = 3.14159265359;

// Adjusted camera settings
const float CAMERA_HEIGHT = 2.0;    
const float CAMERA_DISTANCE = 25.0;  
const float CAMERA_TILT = 0.8;      

// Improved skybox sampling with corrected orientation
vec3 sampleCubemap(vec3 dir) {
    vec3 abs_dir = abs(dir);
    float max_axis = max(max(abs_dir.x, abs_dir.y), abs_dir.z);
    
    vec2 uv;
    vec3 color;
    
    // Scale factor to reduce the "zoomed in" effect
    float scale = 2.0;
    
    // Front/Back
    if (max_axis == abs(dir.z)) {
        if (dir.z > 0.0) {
            uv = vec2(dir.x, dir.y) / abs(dir.z); // Flipped y
            uv = uv * scale + 0.5;
            color = texture(iChannel2, uv).rgb;
        } else {
            uv = vec2(-dir.x, dir.y) / abs(dir.z); // Flipped y
            uv = uv * scale + 0.5;
            color = texture(iChannel1, uv).rgb;
        }
    }
    // Left/Right
    else if (max_axis == abs(dir.x)) {
        if (dir.x > 0.0) {
            uv = vec2(dir.z, dir.y) / abs(dir.x); // Flipped y
            uv = uv * scale + 0.5;
            color = texture(iChannel4, uv).rgb;
        } else {
            uv = vec2(-dir.z, dir.y) / abs(dir.x); // Flipped y
            uv = uv * scale + 0.5;
            color = texture(iChannel3, uv).rgb;
        }
    }
    // Top/Bottom
    else {
        if (dir.y > 0.0) {
            uv = vec2(dir.x, dir.z) / abs(dir.y);
            uv = uv * scale + 0.5;
            color = texture(iChannel5, uv).rgb;
        } else {
            uv = vec2(dir.x, -dir.z) / abs(dir.y);
            uv = uv * scale + 0.5;
            color = texture(iChannel6, uv).rgb;
        }
    }
    
    return color;
}

void mainImage(out vec4 fragColor, in vec2 fragCoord) {
    vec2 uv = (fragCoord - 0.5 * iResolution.xy) / iResolution.y;
    
    // Camera setup with adjusted parameters
    float angle = iTime * 0.5;
    vec3 ro = vec3(
        CAMERA_DISTANCE * sin(angle),
        CAMERA_HEIGHT,
        CAMERA_DISTANCE * cos(angle)
    );
    vec3 ta = vec3(0.0, -2.0, 0.0);
    
    // Camera vectors
    vec3 ww = normalize(ta - ro);
    vec3 uu = normalize(cross(ww, vec3(0.0, 1.0, 0.0)));
    vec3 vv = normalize(cross(uu, ww));
    
    // Ray direction with adjusted field of view
    vec3 rd = normalize(uv.x * uu + uv.y * vv + 1.5 * ww);
    
    // Ground plane intersection
    float t = -(ro.y + 2.0) / rd.y;
    
    vec3 color;
    
    // Check if we hit the ground
    if (t > 0.0 && rd.y < 0.0) {
        vec3 pos = ro + t * rd;
        
        // Adjust checkerboard scale and visibility distance
        vec2 checkUV = pos.xz * 0.05; // Reduced tiling frequency
        
        if (fragCoord.x < iResolution.x * 0.5) {
            // Left side: With MIP mapping
            color = texture(iChannel0, checkUV).rgb;
        } else {
            // Right side: Without MIP mapping (force highest detail)
            color = textureLod(iChannel0, checkUV, 0.0).rgb;
        }
        
        // Improved distance-based shading with closer falloff
        float dist = length(pos.xz);
        float fade = 1.0 - smoothstep(15.0, 30.0, dist); // Reduced fade distance
        
        // Add ambient light to prevent complete darkness
        float ambient = 0.2;
        color = color * max(fade, ambient);
        
        // Add slight fog effect to blend with skybox
        vec3 skyColor = sampleCubemap(rd);
        color = mix(color, skyColor, smoothstep(20.0, 40.0, dist));
    } else {
        // Sample from cubemap skybox
        color = sampleCubemap(rd);
    }
    
    // Add a line to separate the two filtering modes
    if (abs(fragCoord.x - iResolution.x * 0.5) < 1.0) {
        color = vec3(1.0);
    }
    
    fragColor = vec4(color, 1.0);
}