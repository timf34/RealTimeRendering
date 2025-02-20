#iChannel0 "file://checkerboard.png"      // Checkerboard texture
#iChannel1 "file://Skyboxes/back.jpg"     // Skybox back
#iChannel2 "file://Skyboxes/front.jpg"    // Skybox front
#iChannel3 "file://Skyboxes/left.jpg"     // Skybox left
#iChannel4 "file://Skyboxes/right.jpg"    // Skybox right
#iChannel5 "file://Skyboxes/top.jpg"      // Skybox top
#iChannel6 "file://Skyboxes/bottom.jpg"   // Skybox bottom

const float PI = 3.14159265359;

// Function to sample from cubemap
vec3 sampleCubemap(vec3 dir) {
    vec3 abs_dir = abs(dir);
    float max_axis = max(max(abs_dir.x, abs_dir.y), abs_dir.z);
    
    vec2 uv;
    vec3 color;
    
    // Front/Back
    if (max_axis == abs(dir.z)) {
        if (dir.z > 0.0) {
            uv = vec2(dir.x, -dir.y) / abs(dir.z);
            uv = uv * 0.5 + 0.5;
            color = texture(iChannel2, uv).rgb; // front
        } else {
            uv = vec2(-dir.x, -dir.y) / abs(dir.z);
            uv = uv * 0.5 + 0.5;
            color = texture(iChannel1, uv).rgb; // back
        }
    }
    // Left/Right
    else if (max_axis == abs(dir.x)) {
        if (dir.x > 0.0) {
            uv = vec2(dir.z, -dir.y) / abs(dir.x);
            uv = uv * 0.5 + 0.5;
            color = texture(iChannel4, uv).rgb; // right
        } else {
            uv = vec2(-dir.z, -dir.y) / abs(dir.x);
            uv = uv * 0.5 + 0.5;
            color = texture(iChannel3, uv).rgb; // left
        }
    }
    // Top/Bottom
    else {
        if (dir.y > 0.0) {
            uv = vec2(dir.x, dir.z) / abs(dir.y);
            uv = uv * 0.5 + 0.5;
            color = texture(iChannel5, uv).rgb; // top
        } else {
            uv = vec2(dir.x, -dir.z) / abs(dir.y);
            uv = uv * 0.5 + 0.5;
            color = texture(iChannel6, uv).rgb; // bottom
        }
    }
    
    return color;
}

void mainImage(out vec4 fragColor, in vec2 fragCoord) {
    // Normalize coordinates
    vec2 uv = (fragCoord - 0.5 * iResolution.xy) / iResolution.y;
    
    // Camera setup
    vec3 ro = vec3(3.0 * sin(iTime * 0.5), 2.0, 3.0 * cos(iTime * 0.5));
    vec3 ta = vec3(0.0, 0.0, 0.0);
    
    // Camera vectors
    vec3 ww = normalize(ta - ro);
    vec3 uu = normalize(cross(ww, vec3(0.0, 1.0, 0.0)));
    vec3 vv = normalize(cross(uu, ww));
    
    // Ray direction
    vec3 rd = normalize(uv.x * uu + uv.y * vv + 2.0 * ww);
    
    // Ground plane intersection
    float t = -ro.y / rd.y;
    
    vec3 color;
    
    // Check if we hit the ground
    if (t > 0.0 && rd.y < 0.0) {
        // Hit point
        vec3 pos = ro + t * rd;
        
        // Checkerboard UV coordinates
        vec2 checkUV = pos.xz * 0.5;
        
        // Split screen for different filtering modes
        if (fragCoord.x < iResolution.x * 0.5) {
            // Left side: With MIP mapping
            color = texture(iChannel0, checkUV).rgb;
        } else {
            // Right side: Without MIP mapping
            color = textureGrad(iChannel0, checkUV, dFdx(checkUV), dFdy(checkUV)).rgb;
        }
        
        // Add simple shading
        float fade = exp(-length(pos) * 0.1);
        color *= fade;
    } else {
        // Sample from cubemap skybox
        color = sampleCubemap(rd);
    }
    
    // Add a line to separate the two filtering modes
    if (abs(fragCoord.x - iResolution.x * 0.5) < 1.0) {
        color = vec3(1.0);
    }
    
    // Output final color
    fragColor = vec4(color, 1.0);
}