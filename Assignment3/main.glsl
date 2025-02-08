#iChannel0 "file://Skyboxes/front.jpg"     // Using skybox texture as diffuse for now
#iChannel1 "file://Skyboxes/front.jpg"     // Using same texture for normal calculation
#iChannel2 "file://Skyboxes/front.jpg"     // Using same texture for height
#iChannel3 "file://Skyboxes/left.jpg"      // Background skybox

// Adjustable Parameters
const float BUMP_STRENGTH = 0.2;        // Strength of bump mapping (0.0 to 1.0)
const float NORMAL_STRENGTH = 1.0;      // Strength of normal mapping (0.0 to 2.0)
const float SPECULAR_POWER = 32.0;      // Specular highlight sharpness
const float AMBIENT_STRENGTH = 0.1;     // Ambient light intensity

// Animation Parameters
const vec3 ROTATION_SPEEDS = vec3(0.2, 0.3, 0.1);  // Rotation speed for each axis
const float LIGHT_ORBIT_SPEED = 0.5;    // Speed of light movement

// Object Parameters
const float SPHERE_RADIUS = 1.0;
const vec3 LIGHT_COLOR = vec3(1.0, 0.95, 0.8);

// Utility functions for rotation
mat3 rotateX(float angle) {
    float s = sin(angle);
    float c = cos(angle);
    return mat3(
        1, 0, 0,
        0, c, -s,
        0, s, c
    );
}

mat3 rotateY(float angle) {
    float s = sin(angle);
    float c = cos(angle);
    return mat3(
        c, 0, s,
        0, 1, 0,
        -s, 0, c
    );
}

mat3 rotateZ(float angle) {
    float s = sin(angle);
    float c = cos(angle);
    return mat3(
        c, -s, 0,
        s, c, 0,
        0, 0, 1
    );
}

// Sphere SDF
float sdSphere(vec3 p, float r) {
    return length(p) - r;
}

// Scene description
float scene(vec3 p) {
    // Apply rotation animation
    p = rotateY(iTime * ROTATION_SPEEDS.x) * 
        rotateX(iTime * ROTATION_SPEEDS.y) * 
        rotateZ(iTime * ROTATION_SPEEDS.z) * p;
    
    return sdSphere(p, SPHERE_RADIUS);
}

// Calculate normal without bump/normal mapping
vec3 calcBaseNormal(vec3 p) {
    const float eps = 0.001;
    const vec2 h = vec2(eps, 0);
    return normalize(vec3(
        scene(p + h.xyy) - scene(p - h.xyy),
        scene(p + h.yxy) - scene(p - h.yxy),
        scene(p + h.yyx) - scene(p - h.yyx)
    ));
}

// Get UV coordinates for sphere
vec2 getSphereUV(vec3 p) {
    vec3 n = normalize(p);
    float u = 0.5 + atan(n.z, n.x) / (2.0 * 3.14159);
    float v = 0.5 - asin(n.y) / 3.14159;
    return vec2(u, v);
}

// Sample normal map and convert to world space
vec3 getNormalFromMap(vec2 uv, mat3 TBN) {
    vec3 normal = texture(iChannel1, uv).xyz * 2.0 - 1.0;
    normal = normalize(normal * vec3(NORMAL_STRENGTH, NORMAL_STRENGTH, 1.0));
    return normalize(TBN * normal);
}

// Calculate bump mapped normal
vec3 calcBumpedNormal(vec3 p, vec3 normal) {
    vec2 uv = getSphereUV(p);
    
    // Calculate tangent space
    vec3 tangent = normalize(cross(normal, vec3(0.0, 1.0, 0.0)));
    vec3 bitangent = normalize(cross(normal, tangent));
    mat3 TBN = mat3(tangent, bitangent, normal);
    
    // Sample height map for bump mapping
    float height = texture(iChannel2, uv).r;
    
    // Calculate bump mapped normal
    const float bumpScale = BUMP_STRENGTH;
    const float h = 0.01;
    
    vec2 uvDx = uv + vec2(h, 0.0);
    vec2 uvDy = uv + vec2(0.0, h);
    
    float hL = texture(iChannel2, uv - vec2(h, 0.0)).r;
    float hR = texture(iChannel2, uvDx).r;
    float hD = texture(iChannel2, uv - vec2(0.0, h)).r;
    float hU = texture(iChannel2, uvDy).r;
    
    vec3 bumpNormal = normalize(vec3(
        (hL - hR) * bumpScale,
        (hD - hU) * bumpScale,
        2.0
    ));
    
    // Combine bump mapping with normal mapping
    vec3 normalMapNormal = getNormalFromMap(uv, TBN);
    return normalize(mix(normalMapNormal, bumpNormal, 0.5));
}

// Ray marching
float rayMarch(vec3 ro, vec3 rd) {
    float t = 0.0;
    for(int i = 0; i < 100; i++) {
        vec3 p = ro + rd * t;
        float d = scene(p);
        if(d < 0.001 || t > 20.0) break;
        t += d;
    }
    return t;
}

void mainImage(out vec4 fragColor, in vec2 fragCoord) {
    vec2 uv = fragCoord/iResolution.xy;
    vec3 ro = vec3(0.0, 0.0, -4.0);
    vec3 rd = normalize(vec3(((fragCoord - 0.5 * iResolution.xy) / iResolution.y).xy, 1.0));
    
    float t = rayMarch(ro, rd);
    
    if(t < 20.0) {
        vec3 p = ro + rd * t;
        vec3 baseNormal = calcBaseNormal(p);
        vec3 normal = calcBumpedNormal(p, baseNormal);
        
        // Dynamic light position
        vec3 lightPos = vec3(
            3.0 * cos(iTime * LIGHT_ORBIT_SPEED),
            2.0,
            3.0 * sin(iTime * LIGHT_ORBIT_SPEED)
        );
        
        vec3 lightDir = normalize(lightPos - p);
        vec3 viewDir = normalize(ro - p);
        vec3 halfDir = normalize(lightDir + viewDir);
        
        // Get texture color
        vec2 sphereUV = getSphereUV(p);
        vec3 albedo = texture(iChannel0, sphereUV).rgb;
        
        // Calculate lighting
        float diff = max(dot(normal, lightDir), 0.0);
        float spec = pow(max(dot(normal, halfDir), 0.0), SPECULAR_POWER);
        
        // Combine lighting components
        vec3 ambient = AMBIENT_STRENGTH * albedo;
        vec3 diffuse = diff * albedo * LIGHT_COLOR;
        vec3 specular = spec * LIGHT_COLOR;
        
        vec3 color = ambient + diffuse + specular;
        
        // Apply basic tone mapping
        color = color / (color + vec3(1.0));
        
        fragColor = vec4(color, 1.0);
    } else {
        // Background color from skybox
        fragColor = texture(iChannel3, uv);
    }
}